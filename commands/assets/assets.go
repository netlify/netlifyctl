package assets

import (
	"context"
	"fmt"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/spf13/cobra"
)

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "List assets attached to a site",
		Long:  "List assets attached to a site",
	}
	cmd.PersistentFlags().StringP("site-id", "s", "", "site id")

	cmd.AddCommand(setupAddCommand(middlewares))
	cmd.AddCommand(setupInfoCommand(middlewares))

	return middleware.SetupCommand(cmd, listAssets, middlewares)
}

func siteIdForCommand(ctx context.Context, cmd *cobra.Command) (string, error) {
	siteId := cmd.Flag("site-id").Value.String()

	if siteId == "" {
		conf, err := middleware.ChooseSiteConf(ctx, cmd)
		if err != nil {
			return "", err
		}
		siteId = conf.Settings.ID
	}

	if siteId == "" {
		return "", fmt.Errorf("missing site ID to add assets to")
	}

	return siteId, nil
}
