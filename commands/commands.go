package commands

import (
	"github.com/netlify/netlifyctl/commands/assets"
	"github.com/netlify/netlifyctl/commands/deploy"
	initC "github.com/netlify/netlifyctl/commands/init"
	"github.com/netlify/netlifyctl/commands/login"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/commands/sites"
)

func addCommands() {
	loginMiddleware := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.LoggingMiddleware,
		middleware.NoAuthMiddleware,
		middleware.ClientMiddleware,
	}
	middlewares := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
		middleware.ClientMiddleware,
	}

	rootCmd.AddCommand(deploy.Setup(middlewares))
	rootCmd.AddCommand(assets.Setup(middlewares))
	rootCmd.AddCommand(sites.Setup(middlewares))
	rootCmd.AddCommand(login.Setup(loginMiddleware))
	rootCmd.AddCommand(initC.Setup(middlewares))
	rootCmd.AddCommand(versionCmd)
}
