package commands

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/netlify/netlify-go-cli/auth"
	"github.com/netlify/open-api/go/plumbing"
	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List sites in your account",
	Long:  "List sites in your account",
	RunE:  listSites,
}

func listSites(cmd *cobra.Command, args []string) error {
	c := plumbing.NewHTTPClient(nil)
	resp, err := c.Operations.ListSites(nil, auth.ClientCredentials())
	if err != nil {
		return err
	}

	sites := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(sites, "SITE\tURL")
	for _, s := range resp.Payload {
		fmt.Fprintf(sites, "\n%s\t%s", s.Name, s.URL)
	}
	tm.Print(sites)
	tm.Flush()

	return nil
}
