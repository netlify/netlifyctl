package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
	"github.com/netlify/netlifyctl/ui"
	netlify "github.com/netlify/open-api/go/porcelain"
)

type deployCmd struct {
	base          string
	publish       string
	title         string
	functions     string
	siteID        string
	siteName      string
	draft         bool
	preProcessSec int
}

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &deployCmd{}
	ccmd := &cobra.Command{
		Use:     "deploy",
		Aliases: []string{"deploys", "d"},
		Short:   "Deploy your site",
		Long:    "Deploy your site",
	}
	ccmd.Flags().StringVarP(&cmd.base, "base-directory", "b", "", "directory to publish")
	ccmd.Flags().StringVarP(&cmd.publish, "publish-directory", "P", "", "directory to publish")
	ccmd.Flags().StringVarP(&cmd.title, "message", "m", "", "message for the deploy title")
	ccmd.Flags().BoolVarP(&cmd.draft, "draft", "d", false, "draft deploy, not published in production")
	ccmd.Flags().StringVarP(&cmd.functions, "functions", "f", "", "function directory to deploy")
	ccmd.Flags().StringVarP(&cmd.siteID, "site-id", "s", "", "explicitly set a site id instead of relying on configuration")
	ccmd.Flags().StringVarP(&cmd.siteName, "name", "n", "", "search a site by its name instead of relying on configuration")
	ccmd.Flags().IntVarP(&cmd.preProcessSec, "preprocess", "p", 300, "the preprocessing timeout measured in seconds. Default is 5 minutes.")
	return middleware.SetupCommand(ccmd, cmd.deploySite, middlewares)
}

func (dc *deployCmd) deploySite(ctx context.Context, cmd *cobra.Command, args []string) error {
	conf := context.GetSiteConfig(ctx)
	if conf.Settings.ID == "" {
		return errors.New("Failed to load site configuration")
	}

	draft, err := cmd.Flags().GetBool("draft")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'draft'")
	}

	fs, err := cmd.Flags().GetString("functions")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'functions'")
	}
	if fs == "" && conf.Build.Functions != "" {
		fs = conf.Build.Functions
	}

	dir := baseDeploy(cmd, conf)

	if conf.ExistConfFile() && conf.Root() != dir {
		cp, err := conf.CopyConfigFile(dir)
		if err != nil {
			return errors.Wrapf(err, "Error copying netlify.toml to the publish directory: %s", dir)
		}

		if cp != "" {
			defer os.Remove(cp)
		}
	}

	obs := operations.NewDeployObserver()

	client := context.GetClient(ctx)
	options := netlify.DeployOptions{
		Observer:     obs,
		SiteID:       conf.Settings.ID,
		Dir:          dir,
		IsDraft:      draft,
		FunctionsDir: fs,
		Title:        dc.title,
	}
	if dc.preProcessSec > 0 {
		options.PreProcessTimeout = time.Second * time.Duration(dc.preProcessSec)
	}

	logrus.WithFields(logrus.Fields{
		"site":  options.SiteID,
		"path":  options.Dir,
		"draft": options.IsDraft}).Debug("deploying site")

	d, err := client.DeploySite(ctx, options)
	if err != nil {
		return err
	}

	if len(d.Required) > 0 {
		ready, err := client.WaitUntilDeployReady(ctx, d, options.PreProcessTimeout)
		if err != nil {
			return err
		}
		d = ready
	}

	obs.Finish()

	u := d.SslURL
	if d.Context != "production" {
		u = d.DeploySslURL
	}
	fmt.Printf("Deploy done  %s\n", ui.WorldCheck())
	ui.Bold("    %s\n", u)

	return nil
}

func baseDeploy(cmd *cobra.Command, conf *configuration.Configuration) string {
	bd, err := cmd.Flags().GetString("base-directory")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'base-directory'")
	}

	if bd != "" {
		ui.Warning("the base-directory behavior has been deprecated and it will be changed in 0.5.0 version, use the publish-directory flag instead")
		return bd
	}

	pd, err := cmd.Flags().GetString("publish-directory")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get string flag: 'publish-directory'")
	}

	if pd != "" {
		return pd
	}

	s := conf.Settings
	path := s.Path

	if path == "" && conf.Build.Publish != "" {
		path = conf.Build.Publish
	}

	if path == "" {
		path, err = ui.AskForInput("What path would you like deployed?", ".")
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get deploy path")
		}

		logrus.Debugf("Got new path from the user %s", s.Path)
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(conf.Root(), path)
		logrus.Debugf("Relative path detected, going to deploy: '%s'", path)
	}

	s.Path = path

	return path
}
