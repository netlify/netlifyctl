package forms

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
		Use:     "form",
		Aliases: []string{"forms", "f"},
		Short:   "List forms",
		Long:    "List forms",
	}

	cmd.AddCommand(setupFormSubmissions(middlewares))

	return middleware.SetupCommand(cmd, listForms, middlewares)
}

func listForms(ctx context.Context, cmd *cobra.Command, args []string) error {
	client := context.GetClient(ctx)

	forms, err := client.ListForms(ctx, nil)
	if err != nil {
		return err
	}

	t := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "Form Name\tForm ID\tSite ID\tSubmissions\n")
	for _, f := range forms {
		fmt.Fprintf(buffer, "%s\t%s\t%s\t%d\n", f.Name, f.ID, f.SiteID, f.SubmissionCount)
	}
	buffer.WriteTo(t)
	t.Flush()

	return nil
}
