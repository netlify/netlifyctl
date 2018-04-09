package init

import (
	"fmt"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
	"github.com/spf13/cobra"
)

var manual bool

func Setup(middlewares []middleware.Middleware) *cobra.Command {
	ccmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"cd", "i"},
		Short:   "Configure continuous deployment",
		Long:    "Configure continuous deployment",
	}
	ccmd.Flags().BoolVarP(&manual, "manual", "m", false, "Step by step setup (no Git Provider permissions required)")

	return middleware.SetupCommand(ccmd, initSite, middlewares)
}

func initSite(ctx context.Context, cmd *cobra.Command, args []string) error {
	host, err := getRepo()
	if err != nil {
		return err
	}

	tt := ui.NewTaskTracker()
	var c configurator
	if manual {
		c = manualConfigurator{host}
		tt = ui.NewTaskTrackerWithTerm(false)
	} else {
		ec, err := loadConfigurator(ctx, host.URL)
		if err != nil {
			return err
		}
		c = ec
	}

	site, err := operations.ChooseOrCreateSite(ctx, cmd)
	if err != nil {
		return err
	}
	if site == nil {
		return nil
	}

	dir, err := ui.AskForInput("Directory to deploy (blank for current dir):", ".")
	if err != nil {
		return err
	}
	buildCmd, err := ui.AskForInput("Your build command (hugo build/yarn run build/etc):", "")
	if err != nil {
		return err
	}

	info, err := c.RepoInfo(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("\nConfiguration:\n\n")
	fmt.Printf("    Repository: %s\n", host.Remote)
	fmt.Printf("    Production branch: %s\n", info.RepoBranch)
	fmt.Printf("    Publishing directory: %s\n", dir)
	fmt.Printf("    Build command: %s\n\n", buildCmd)

	if !ui.AskForConfirmation("Continue?") {
		return nil
	}

	err = ui.TrackWithTracker("Configuring Continuous Deployment ... ", "Success! Whenever you push to git, Netlify will build and deploy your site", tt, func() error {
		client := context.GetClient(ctx)
		key, err := client.CreateDeployKey(ctx)
		if err != nil {
			return err
		}

		if err := c.SetupDeployKey(ctx, key); err != nil {
			return err
		}

		info.DeployKeyID = key.ID
		if dir != "" {
			info.Dir = dir
		}
		if buildCmd != "" {
			info.Cmd = buildCmd
		}

		setup := &models.SiteSetup{
			Site: *site,
			Repo: info,
		}
		uSite, err := client.UpdateSite(ctx, setup)
		if err != nil {
			return err
		}

		return c.SetupWebHook(ctx, uSite)
	})
	if err != nil {
		return err
	}

	ui.Bold("    %s\n", site.SslURL)

	return nil
}
