package sites

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/spf13/cobra"
)

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "site",
		Aliases: []string{"sites", "s"},
		Short:   "Handle site operations",
		Long:    "Handle site operations",
	}

	cmd.AddCommand(setupUpdateCommand(middlewares))
	cmd.AddCommand(setupCreateCommand(middlewares))

	return middleware.SetupCommand(cmd, listSites, middlewares)
}

func listSites(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)

	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		return err
	}

	t := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "Site\tUrl\n")
	for _, s := range sites {
		fmt.Fprintf(buffer, "%s\t%s\n", s.Name, s.URL)
	}
	buffer.WriteTo(t)
	t.Flush()

	return nil
}
