package init

import (
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
)

type gitlabConfigurator struct {
	accessToken string
}

func (c *gitlabConfigurator) Login() error {
	return nil
}

func (c *gitlabConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	return nil
}

func (c *gitlabConfigurator) SetupWebHook(ctx context.Context, site *models.Site) error {
	return nil
}

func (c *gitlabConfigurator) RepoInfo(ctx context.Context) (*models.RepoSetup, error) {
	return nil, nil
}
