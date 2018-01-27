package middleware

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/netlify/open-api/go/porcelain"
	apiContext "github.com/netlify/open-api/go/porcelain/context"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
)

const (
	defaultAPIPath = "/api/v1"
	debugLogFile   = "netlifyctl-debug.log"
)

type CommandFunc func(context.Context, *cobra.Command, []string) error
type Middleware func(CommandFunc) CommandFunc

func SetupCommand(cmd *cobra.Command, f CommandFunc, m []Middleware) *cobra.Command {
	cmd.RunE = NewRunFunc(f, m)
	return cmd
}

func NewRunFunc(f CommandFunc, mm []Middleware) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		runf := f
		for _, m := range mm {
			runf = m(runf)
		}

		return runf(ctx, cmd, args)
	}
}

func DebugMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		b := new(bytes.Buffer)
		// Enable open-api debug mode and disable it after running the command.
		os.Setenv("DEBUG", "1")
		defer os.Unsetenv("DEBUG")

		// Enable debug logging
		logrus.SetOutput(b)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.WithFields(logrus.Fields{"command": c.Use, "arguments": args}).Debug("PreRun")

		logrus.Debug("configure debug middleware")

		dump, err := c.Root().Flags().GetBool("debug")
		if err != nil {
			return err
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			if dump {
				dumpDebug(b)
				fmt.Printf("\nDebug log dumped to %s\n", debugLogFile)
				fmt.Println("This log includes full recordings of HTTP requests with credentials, be careful if you share it")
			}
			os.Exit(0)
		}()

		// Run command
		if err := cmd(ctx, c, args); err != nil {
			logrus.WithError(err).Error("command failed")
			if err := dumpDebug(b); err != nil {
				return err
			}
			return fmt.Errorf("There was an error running this command.\nDebug log dumped to %s\nThis log includes full recordings of HTTP requests with credentials, be careful if you share it", debugLogFile)
		}

		if dump {
			if err := dumpDebug(b); err != nil {
				return err
			}
			fmt.Printf("Debug log dumped to %s\n", debugLogFile)
			fmt.Println("This log includes full recordings of HTTP requests with credentials, be careful if you share it")
		}

		return nil
	}
}

func LoggingMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		logrus.Debug("configure logging middleware")
		entry := logrus.NewEntry(logrus.StandardLogger())
		entry.Debugf("setup logger middleware: %v", entry.Level)
		ctx = apiContext.WithLogger(ctx, entry)

		return cmd(ctx, c, args)
	}
}

func AuthMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		logrus.Debug("configure auth middleware")
		creds := auth.ClientCredentials()
		logrus.WithField("credentials", creds).Debug("setup credentials")

		ctx = apiContext.WithAuthInfo(ctx, creds)

		return cmd(ctx, c, args)
	}
}

func NoAuthMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		logrus.Debug("configure no auth middleware")
		creds := auth.NoCredentials()
		logrus.WithField("credentials", creds).Debug("setup credentials")

		ctx = apiContext.WithAuthInfo(ctx, creds)

		return cmd(ctx, c, args)
	}
}

func ClientMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		logrus.Debug("configure client middleware")
		var transport *apiClient.Runtime

		if endpoint := c.Flag("endpoint"); endpoint != nil {
			if v := endpoint.Value.String(); v != "" {
				logrus.WithField("endpoint", v).Debug("setup API endpoint")

				u, err := url.Parse(v)
				if err != nil {
					return err
				}

				if u.Scheme == "" {
					u.Scheme = "http"
				}

				if u.Path == "" {
					u.Path = defaultAPIPath
				}

				transport = apiClient.NewWithClient(u.Host, u.Path, []string{u.Scheme}, httpClient())
			}
		}

		if transport == nil {
			logrus.WithField("endpoint", "https://api.netlify.com").Debug("setup default API endpoint")

			transport = apiClient.NewWithClient("api.netlify.com", "", []string{"https"}, httpClient())
		}

		logger := apiContext.GetLogger(ctx)
		transport.SetDebug(true)
		transport.SetLogger(logger)

		client := porcelain.New(transport, strfmt.Default)
		ctx = context.WithClient(ctx, client)

		return cmd(ctx, c, args)
	}
}

func SiteConfigMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		logrus.Debug("configure site middleware")
		var siteId string
		if siteIdFlag := c.Flag("site-id"); siteIdFlag != nil {
			siteId = siteIdFlag.Value.String()
		}

		configFile := c.Root().Flag("config").Value.String()
		conf, err := configuration.Load(configFile)
		if err != nil {
			return err
		}
		if siteId != "" && conf.Settings.ID != siteId {
			conf.Settings.ID = siteId
		}

		if conf.Settings.ID == "" {
			logrus.Debug("Querying for existing sites")
			// we don't know the site - time to try and get its id
			site, err := operations.ChooseOrCreateSite(ctx, c)

			// Ensure that the site ID is always saved,
			// even when there is a provision error.
			if site != nil {
				conf.Settings.ID = site.ID
				configuration.Save(configFile, conf)
			}

			if err != nil {
				return err
			}
		}
		ctx = context.WithSiteConfig(ctx, conf)

		return cmd(ctx, c, args)
	}
}

func httpClient() *http.Client {
	protoUpgrade := map[string]func(string, *tls.Conn) http.RoundTripper{
		"ignore-h2": func(string, *tls.Conn) http.RoundTripper { return nil },
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSNextProto:          protoUpgrade,
	}

	return &http.Client{Transport: tr}
}

func dumpDebug(b *bytes.Buffer) error {
	return ioutil.WriteFile(debugLogFile, b.Bytes(), 0644)
}
