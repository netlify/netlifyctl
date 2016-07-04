package middleware

import (
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/Sirupsen/logrus"
	logContext "github.com/docker/distribution/context"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/porcelain"
	authContext "github.com/netlify/open-api/go/porcelain/context"
	"github.com/spf13/cobra"
)

const defaultAPIPath = "/api/v1"

type CommandFunc func(context.Context, *cobra.Command, []string) error
type Middleware func(CommandFunc) CommandFunc

func NewRunFunc(f CommandFunc, mm []Middleware) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := logContext.Background()

		runf := f
		for _, m := range mm {
			runf = m(runf)
		}

		return runf(ctx, cmd, args)
	}
}

func LoggingMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		entry := logrus.NewEntry(logrus.StandardLogger())
		logrus.WithField("log_level", "debug").Debug("setup logger middleware")

		ctx = logContext.WithLogger(ctx, entry)

		return cmd(ctx, c, args)
	}
}

func AuthMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		creds := auth.ClientCredentials()
		logrus.WithField("credentials", creds).Debug("setup credentials")

		ctx = authContext.WithAuthInfo(ctx, creds)

		return cmd(ctx, c, args)
	}
}

func ClientMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		var client *porcelain.Netlify

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

				transport := httptransport.New(u.Host, u.Path, []string{u.Scheme})
				client = porcelain.New(transport, strfmt.Default)
			}
		}

		if client == nil {
			logrus.WithField("endpoint", "https://api.netlify.com").Debug("setup default API endpoint")
			client = porcelain.NewHTTPClient(nil)
		}

		ctx = context.WithClient(ctx, client)

		return cmd(ctx, c, args)
	}
}
