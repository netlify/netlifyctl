package assets

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

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

	t := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "Id\tName\tVisibility\tUrl")
	for _, a := range assets {
		fmt.Fprintf(buffer, "\n%s\t%s\t%s\t%s", a.ID, a.Name, a.Visibility, a.URL)
	}
	buffer.WriteTo(t)
	t.Flush()

	return nil
}
