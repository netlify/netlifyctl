package deploy

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"path/filepath"
	"strings"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
	netlify "github.com/netlify/open-api/go/porcelain"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	base  string
	title string
	draft bool
}

func Setup() (*cobra.Command, middleware.CommandFunc) {
	cmd := &deployCmd{}
	ccmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your site",
		Long:  "Deploy your site",
	}
	ccmd.Flags().StringVarP(&cmd.base, "base-directory", "b", "", "directory to publish")
	ccmd.Flags().StringVarP(&cmd.title, "message", "m", "", "message for the deploy title")
	ccmd.Flags().BoolVarP(&cmd.draft, "draft", "d", false, "draft deploy, not published in production")

	return ccmd, cmd.deploySite
}

func (dc *deployCmd) deploySite(ctx context.Context, cmd *cobra.Command, args []string) error {
	var configFile = cmd.Root().Flag("config").Value.String()
	var conf, err = configuration.Load(configFile)
	if err != nil {
		return err
	}
	client := context.GetClient(ctx)

	if conf.Settings.ID == "" {

		logrus.Debug("Querying for existing sites")
		// we don't know the site - time to try and get its id
		site, err := operations.ChooseOrCreateSite(client, ctx)

		// Ensure that the site ID is always saved,
		// even when there is a provision error.
		if site != nil {
			conf.Settings.ID = site.ID
			configuration.Save(configFile, conf)
		}
		if err != nil {
			return err
		}

		fmt.Printf("=>  Domain ready, deploying assets to %s\n", site.Name)
	}

	options := netlify.DeployOptions{
		SiteID: conf.Settings.ID,
		Dir:    baseDeploy(cmd, conf),
	}

	configuration.Save(configFile, conf)

	draft, err := cmd.Flags().GetBool("draft")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'draft'")
	}
	options.IsDraft = draft

	logrus.WithFields(logrus.Fields{
		"site":  options.SiteID,
		"path":  options.Dir,
		"draft": options.IsDraft}).Debug("deploying site")

	d, err := client.DeploySite(ctx, options)
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
	var path = s.Path
	if path == "" {
		path, err = operations.AskForInput("What path would you like deployed?", ".")
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get deploy path")
		}

		s.Path = path
		logrus.Debugf("Got new path from the user %s", s.Path)
	}

	if !strings.HasPrefix(s.Path, "/") {
		path = filepath.Join(conf.Root(), s.Path)
		logrus.Debugf("Relative path detected, going to deploy: '%s'", path)
	}

	return path
}
