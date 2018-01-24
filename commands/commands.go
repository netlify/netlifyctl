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
	middlewares := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
		middleware.ClientMiddleware,
	}

	loginMiddlewares := []middleware.Middleware{
		middleware.DebugMiddleware,
		middleware.LoggingMiddleware,
		middleware.NoAuthMiddleware,
		middleware.ClientMiddleware,
	}

	siteMiddlewares := append(middlewares, middleware.SiteConfigMiddleware)

	rootCmd.AddCommand(deploy.Setup(siteMiddlewares))
	rootCmd.AddCommand(assets.Setup(siteMiddlewares))
	rootCmd.AddCommand(sites.Setup(middlewares))
	rootCmd.AddCommand(initC.Setup(middlewares))
	rootCmd.AddCommand(login.Setup(loginMiddlewares))
	rootCmd.AddCommand(versionCmd)
}
