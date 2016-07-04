package sites

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/spf13/cobra"
)

func Setup() (*cobra.Command, middleware.CommandFunc) {
	return &cobra.Command{
		Use:   "sites",
		Short: "List sites in your account",
		Long:  "List sites in your account",
	}, listSites
}

func listSites(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)

	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		return err
	}

	table := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(table, "SITE\tURL")
	for _, s := range sites {
		fmt.Fprintf(table, "\n%s\t%s", s.Name, s.URL)
	}
	tm.Print(table)
	tm.Flush()

	return nil
}
