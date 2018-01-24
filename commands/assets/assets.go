package assets

import (
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/spf13/cobra"
)

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assets",
		Aliases: []string{"asset", "a"},
		Short:   "List assets attached to a site",
		Long:    "List assets attached to a site",
	}
	cmd.PersistentFlags().StringP("site-id", "s", "", "site id")

	cmd.AddCommand(setupAddCommand(middlewares))
	cmd.AddCommand(setupInfoCommand(middlewares))

	return middleware.SetupCommand(cmd, listAssets, middlewares)
}
