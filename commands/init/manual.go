package init

import (
	"fmt"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
)

type manualConfigurator struct {
	gitProvider *gitRepo
}

func (c manualConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	fmt.Println("Give this Netlify SSH public key access to your repo:\n\n")
	fmt.Printf("%s\n\n", deployKey.PublicKey)

	ui.AskForConfirmation("Continue?")
	return nil
}

func (c manualConfigurator) SetupWebHook(ctx context.Context, site *models.Site) error {
	fmt.Printf("Configure the following webhook for your repository:\n\n")
	fmt.Printf("%s\n\n", site.WebHook)

	ui.AskForConfirmation("Continue?")
	fmt.Println("Success! Whenever you push to git, Netlify will build and deploy your site")
	fmt.Printf("  %s", site.URL)
	return nil
}

func (c manualConfigurator) RepoInfo(ctx context.Context) (*models.RepoSetup, error) {
	branch := c.gitProvider.CurrentBranch

	return &models.RepoSetup{
		Provider:        "manual",
		Repo:            c.gitProvider.Remote,
		Branch:          branch,
		AllowedBranches: []string{branch},
	}, nil
	return nil, nil
}
