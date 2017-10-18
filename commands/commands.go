package commands

import (
	"github.com/netlify/netlifyctl/commands/assets"
	"github.com/netlify/netlifyctl/commands/deploy"
	"github.com/netlify/netlifyctl/commands/login"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/commands/sites"
)

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
	rootCmd.AddCommand(middleware.SetupCommand(dCmd, dFunc, middlewares))

	rootCmd.AddCommand(assets.Setup(middlewares))
	rootCmd.AddCommand(sites.Setup(middlewares))
	rootCmd.AddCommand(login.Setup(loginMiddleware))
	rootCmd.AddCommand(versionCmd)
}
