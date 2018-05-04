package commands

import (
	"github.com/netlify/netlifyctl/commands/assets"
	"github.com/netlify/netlifyctl/commands/deploy"
	"github.com/netlify/netlifyctl/commands/forms"
	initC "github.com/netlify/netlifyctl/commands/init"
	"github.com/netlify/netlifyctl/commands/login"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/commands/sites"
)

func addCommands() {
	middlewares := []middleware.Middleware{
		middleware.ClientMiddleware,
		middleware.AuthMiddleware,
		middleware.LoggingMiddleware,
		middleware.DebugMiddleware,
	}

	loginMiddlewares := []middleware.Middleware{
		middleware.ClientMiddleware,
		middleware.NoAuthMiddleware,
		middleware.LoggingMiddleware,
		middleware.DebugMiddleware,
	}

	siteMiddlewares := append([]middleware.Middleware{middleware.SiteConfigMiddleware}, middlewares...)

	rootCmd.AddCommand(deploy.Setup(siteMiddlewares))
	rootCmd.AddCommand(assets.Setup(siteMiddlewares))
	rootCmd.AddCommand(sites.Setup(middlewares))
	rootCmd.AddCommand(forms.Setup(middlewares))
	rootCmd.AddCommand(initC.Setup(middlewares))
	rootCmd.AddCommand(login.Setup(loginMiddlewares))
	rootCmd.AddCommand(versionCmd)
}
