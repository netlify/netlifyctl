package assets

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/spf13/cobra"
)

func listAssets(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)

	siteId, err := siteIdForCommand(ctx, cmd)
	if err != nil {
		return err
	}

	params := operations.NewListSiteAssetsParams().WithSiteID(siteId)
	assets, err := client.ListSiteAssets(ctx, params)
	if err != nil {
		return err
	}

	table := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(table, "ID\tNAME\tVISIBILITY\tURL")
	for _, a := range assets {
		fmt.Fprintf(table, "\n%s\t%s\t%s\t%s", a.ID, a.Name, a.Visibility, a.URL)
	}
	tm.Print(table)
	tm.Flush()

	return nil
}
