package streaming

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/errors"
)

func Setup() (*cobra.Command, middleware.CommandFunc) {
	ccmd := &cobra.Command{
		Use:   "stream <deploy_id>",
		Short: "stream a deploy log",
		Long:  "stream your deploy log",
	}

	return ccmd, streamLogs
}

func streamLogs(ctx context.Context, cmd *cobra.Command, args []string) error {
	// extract the deployID
	if len(args) != 1 {
		return errors.ArgumentErrorF("Must provide a deployID")
	}

	deployID := args[0]

	client := context.GetClient(ctx)
	msgs, shutdown, err := client.StreamBuildLog(ctx, deployID)
	if err != nil {
		return err
	}

	// TODO traps to shutdown
	_ = shutdown

	for m := range msgs {
		var msg string
		if m.ErrorMessage != "" {
			msg = "ERROR: " + m.ErrorMessage
		} else if m.Completed {
			msg = "COMPLETED"
		} else {
			if m.Phase != "" {
				msg += fmt.Sprintf("%s - ", m.Phase)
			}
			msg += m.Message
		}

		fmt.Printf("%s: %s : %s\n", m.Time, m.DeployID, msg)
	}
	return nil
}
