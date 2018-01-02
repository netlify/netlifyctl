package init

import (
	"fmt"
	"os"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
)

type manualConfigurator struct {
	gitProvider *gitRepo
}

func (c manualConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	fmt.Println("\nGive this Netlify SSH public key access to your repository:\n\n")
	ui.Bold("%s\n\n", deployKey.PublicKey)

	if !ui.AskForConfirmation("Continue?") {
		os.Exit(0)
	}
	return nil
}

func (c manualConfigurator) SetupWebHook(ctx context.Context, site *models.Site) error {
	fmt.Printf("\nConfigure the following webhook for your repository:\n\n")
	ui.Bold("    %s\n\n", site.DeployHook)

	if !ui.AskForConfirmation("Continue?") {
		os.Exit(0)
	}

	return nil
}

func (c manualConfigurator) RepoInfo(ctx context.Context) (*models.RepoInfo, error) {
	branch := c.gitProvider.CurrentBranch

	return &models.RepoInfo{
		Provider:        "manual",
		RepoPath:        c.gitProvider.Remote,
		RepoBranch:      branch,
		AllowedBranches: []string{branch},
	}, nil
	return nil, nil
}
