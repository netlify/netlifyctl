package init

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
	"github.com/pkg/errors"
)

type manualConfigurator struct {
	gitProvider *gitRepo
}

func (c manualConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	fmt.Printf("\nGive this Netlify SSH public key access to your repository:\n\n")
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

	repoPath := c.gitProvider.Remote
	if strings.HasPrefix(repoPath, "https://") {
		p, err := url.Parse(repoPath)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid repository url: %s", repoPath)
		}
		hp := strings.SplitN(p.Host, ":", 2)
		repoPath = fmt.Sprintf("git@%s:%s", hp[0], strings.TrimPrefix(p.Path, "/"))
		if !strings.HasSuffix(repoPath, ".git") {
			repoPath = fmt.Sprintf("%s.git", repoPath)
		}
	}

	return &models.RepoInfo{
		Provider:        "manual",
		RepoPath:        repoPath,
		RepoBranch:      branch,
		AllowedBranches: []string{branch},
	}, nil
}
