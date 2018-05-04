package forms

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type submissionsCmd struct {
	formID string
}

func setupFormSubmissions(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &submissionsCmd{}
	ccmd := &cobra.Command{
		Use:   "submissions <-i FORM_ID> ...",
		Short: "list form submissions",
		Long:  "list form submissions",
	}
	ccmd.Flags().StringVarP(&cmd.formID, "form-id", "i", "", "form's id")

	return middleware.SetupCommand(ccmd, cmd.listFormSubmissions, middlewares)
}

func (sc *submissionsCmd) listFormSubmissions(ctx context.Context, cmd *cobra.Command, args []string) error {
	formID, err := cmd.Flags().GetString("form-id")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'form-id'")
	}

	client := context.GetClient(ctx)

	submissions, err := client.ListFormSubmissions(ctx, formID)

	if err != nil {
		return err
	}

	t := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "Timestamp\tName\tBody\n")
	for _, f := range submissions {
		fmt.Fprintf(buffer, "%s\t%s\t%s\n", f.CreatedAt, f.Name, f.Body)
	}
	buffer.WriteTo(t)
	t.Flush()

	return nil
}
