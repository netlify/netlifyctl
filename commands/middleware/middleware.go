package middleware

import (
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/Sirupsen/logrus"
	logContext "github.com/docker/distribution/context"
	strfmt "github.com/go-openapi/strfmt"
	authContext "github.com/netlify/open-api/go/porcelain/context"
	"github.com/spf13/cobra"

	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/porcelain"
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

		ctx = logContext.WithLogger(ctx, entry)

		logrus.WithField("log_level", "debug").Debug("setup logger middleware")
		return cmd(ctx, c, args)
	}
}

func AuthMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		ctx = authContext.WithAuthToken(ctx, auth.AuthToken())

		logrus.Debug("setup credentials")
		return cmd(ctx, c, args)
	}
}

func UserAgentMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		ctx = authContext.WithUserAgent(ctx, "netlifyctl")

		logrus.Debug("setup user agent")
		return cmd(ctx, c, args)
	}
}

func ClientMiddleware(cmd CommandFunc) CommandFunc {
	return func(ctx context.Context, c *cobra.Command, args []string) error {
		var client *porcelain.Netlify

		if endpoint, err := c.PersistentFlags().GetString("endpoint"); err != nil && endpoint != "" {
			logrus.WithField("endpoint", endpoint).Debug("setup API endpoint")

			u, err := url.Parse(endpoint)
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

		if client == nil {
			logrus.WithField("endpoint", "https://api.netlify.com").Debug("setup default API endpoint")
			client = porcelain.NewHTTPClient(nil)
		}

		if streamEndpoint, err := c.PersistentFlags().GetString("streaming"); err != nil && streamEndpoint != "" {
			err := client.SetStreamingEndpoint(streamEndpoint)
			if err != nil {
				return err
			}
		}

		ctx = context.WithClient(ctx, client)

		return cmd(ctx, c, args)
	}
}
