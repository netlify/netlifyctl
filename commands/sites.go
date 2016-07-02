package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List sites in your account",
	Long:  "List sites in your account",
	RunE:  listSites,
}

func listSites(cmd *cobra.Command, args []string) error {
	c := porcelain.NewHTTPClient(nil)
	ctx := auth.NewContext()

	sites, err := c.ListSites(ctx, nil)
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
