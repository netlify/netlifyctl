package commands

import (
	"github.com/netlify/netlifyctl/commands/deploy"
	"github.com/netlify/netlifyctl/commands/login"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/commands/sites"
	"github.com/spf13/cobra"
)

func setupRunE(cmd *cobra.Command, f middleware.CommandFunc, m []middleware.Middleware) *cobra.Command {
	cmd.RunE = middleware.NewRunFunc(f, m)
	return cmd
}

func addCommands() {
	loginMiddleware := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.NoAuthMiddleware,
		middleware.ClientMiddleware,
		middleware.LoggingMiddleware,
	}
	middlewares := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.AuthMiddleware,
		middleware.ClientMiddleware,
		middleware.LoggingMiddleware,
	}

	dCmd, dFunc := deploy.Setup()
	rootCmd.AddCommand(setupRunE(dCmd, dFunc, middlewares))

	rootCmd.AddCommand(sites.Setup(middlewares))
	rootCmd.AddCommand(login.Setup(loginMiddleware))
	rootCmd.AddCommand(versionCmd)
}
