package assets

import (
	"encoding/json"
	"fmt"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/spf13/cobra"
)

type assetsShowCmd struct {
	withSignature bool
}

func setupInfoCommand(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &assetsShowCmd{}
	ccmd := &cobra.Command{
		Use:   "info [ASSET ID 1] [ASSET ID 2] ...",
		Short: "Show information for an asset or a group of them",
		Long:  "Show information for an asset or a group of them",
	}
	ccmd.Flags().BoolVarP(&cmd.withSignature, "with-public-signature", "p", false, "show public signature of private assets")

	return middleware.SetupCommand(ccmd, cmd.showAssets, middlewares)
}

func (c *assetsShowCmd) showAssets(ctx context.Context, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing asset ids to show")
	}

	siteId, err := siteIdForCommand(ctx, cmd)
	if err != nil {
		return err
	}

	client := context.GetClient(ctx)
	for i, arg := range args {
		params := operations.NewGetSiteAssetInfoParams().WithSiteID(siteId).WithAssetID(arg)

		asset, err := client.ShowSiteAssetInfo(ctx, params, c.withSignature)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(asset, "", "\t")
		if err != nil {
			return err
		}

		fmt.Print(string(b))
		if i+1 < len(args) {
			fmt.Print("\n")
		}
	}
	return nil
}
