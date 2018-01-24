package assets

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/spf13/cobra"
)

func listAssets(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)
	conf := context.GetSiteConfig(ctx)
	if conf.Settings.ID == "" {
		return errors.New("Failed to load site configuration")
	}

	params := operations.NewListSiteAssetsParams().WithSiteID(conf.Settings.ID)
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
