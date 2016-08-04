package commands

import (
	"github.com/netlify/netlifyctl/commands/deploy"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/commands/sites"
	"github.com/spf13/cobra"
)

func setupRunE(cmd *cobra.Command, f middleware.CommandFunc, m []middleware.Middleware) *cobra.Command {
	cmd.RunE = middleware.NewRunFunc(f, m)
	return cmd
}

func addCommands() {
	middlewares := []middleware.Middleware{
		middleware.AuthMiddleware,
		middleware.UserAgentMiddleware,
		middleware.ClientMiddleware,
		middleware.LoggingMiddleware,
	}

	sCmd, sFunc := sites.Setup()
	rootCmd.AddCommand(setupRunE(sCmd, sFunc, middlewares))

	dCmd, dFunc := deploy.Setup()
	rootCmd.AddCommand(setupRunE(dCmd, dFunc, middlewares))
}

// ArgumentError is for when there is a problem with the args passed in the CLI
type ArgumentError struct {
	errmMsg string
}

func (e ArgumentError) Error() string {
	return e.errmMsg
}
