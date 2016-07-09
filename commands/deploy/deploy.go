package deploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	base string
}

func Setup() (*cobra.Command, middleware.CommandFunc) {
	cmd := &deployCmd{}
	ccmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your site",
		Long:  "Deploy your site",
	}
	ccmd.Flags().StringVarP(&cmd.base, "base-directory", "b", "", "directory to publish")

	return ccmd, cmd.deploySite
}

func (*deployCmd) deploySite(ctx context.Context, cmd *cobra.Command, args []string) error {
	conf, err := configuration.Load()
	if err != nil {
		return err
	}
	client := context.GetClient(ctx)

	if conf.Settings.ID == "" && operations.ConfirmCreateSite(cmd) {
		newSite, err := operations.CreateSite(cmd, client, ctx)
		// Ensure that the site ID is always saved,
		// even when there is a provision error.
		if newSite != nil {
			conf.Settings.ID = newSite.ID
			configuration.Save(conf)
		}

		if err != nil {
			return err
		}

		fmt.Println("=> Domain ready, deploying assets now")
	}

	s := conf.Settings
	id := s.ID

	base := baseDeploy(cmd, s)
	logrus.WithFields(logrus.Fields{"site": id, "root": base}).Debug("deploy site")

	d, err := client.DeploySite(ctx, id, base)
	if err != nil {
		return err
	}

	ready, err := client.WaitUntilDeployReady(ctx, d)
	if err != nil {
		return err
	}
	fmt.Printf("=> Done, your website is live in %s\n", ready.URL)

	return nil
}

func baseDeploy(cmd *cobra.Command, conf configuration.Settings) string {
	f := cmd.Flag("base-directory")
	if f == nil || f.Value.String() == "" {
		return conf.Root()
	}
	return f.Value.String()
}
