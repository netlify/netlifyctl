package deploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"path/filepath"
	"strings"

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

	id := conf.Settings.ID

	path := baseDeploy(cmd, conf)
	configuration.Save(conf)
	logrus.WithFields(logrus.Fields{"site": id, "path": path}).Debug("deploying site")

	d, err := client.DeploySite(ctx, id, path)
	if err != nil {
		return err
	}

	if len(d.Required) > 0 {
		ready, err := client.WaitUntilDeployReady(ctx, d)
		if err != nil {
			return err
		}
		d = ready
	}
	fmt.Printf("=> Done, your website is live in %s\n", d.URL)

	return nil
}

func baseDeploy(cmd *cobra.Command, conf *configuration.Configuration) string {
	bd, err := cmd.Flags().GetString("base-directory")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'base-directory'")
	}

	if bd != "" {
		return bd
	}
	s := conf.Settings
	path := s.Path
	if path == "" {
		path, err = operations.AskForInput("What path would you like deployed?", ".")
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get deploy path")
		}

		s.Path = path
		logrus.Debugf("Got new path from the user %s", s.Path)
	}

	if !strings.HasPrefix(s.Path, "/") {
		path := filepath.Join(conf.Root(), s.Path)
		logrus.Debugf("Relative path detected, going to deploy: '%s'", path)
	}

	return path
}
