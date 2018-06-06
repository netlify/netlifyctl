package init

import (
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	hub "github.com/github/hub/github"
	"github.com/google/go-github/github"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
)

const (
	deployKeyTitle = "Netlify Deploy Key"
)

type githubConfigurator struct {
	client *github.Client
	owner  string
	repo   string
}

func newGitHubConfigurator(ctx context.Context, url *url.URL) (*githubConfigurator, error) {
	h, err := hub.CurrentConfig().PromptForHost(githubHost)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: h.AccessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return newGitHubConfiguratorWithClient(ctx, url, client)
}

func newGitHubConfiguratorWithClient(ctx context.Context, url *url.URL, client *github.Client) (*githubConfigurator, error) {
	parts := strings.SplitN(strings.TrimLeft(url.Path, "/"), "/", 2)

	return &githubConfigurator{
		client: client,
		owner:  parts[0],
		repo:   parts[1],
	}, nil
}

func (c *githubConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	key := &github.Key{
		Title: github.String(deployKeyTitle),
		Key:   github.String(deployKey.PublicKey),
	}
	_, _, err := c.client.Repositories.CreateKey(ctx, c.owner, c.repo, key)
	return err
}

func (c *githubConfigurator) SetupWebHook(ctx context.Context, site *models.Site) error {
	hooks, _, err := c.client.Repositories.ListHooks(ctx, c.owner, c.repo, nil)
	if err != nil {
		return err
	}

	// Do not try to install the webhook twice.
	if hooks != nil && len(hooks) > 0 {
		for _, h := range hooks {
			if h.Config["url"] == site.DeployHook {
				return nil
			}
		}
	}

	hook := &github.Hook{
		Name:   github.String("web"),
		Events: []string{"push", "pull_request", "delete"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          site.DeployHook,
			"content_type": "json",
		},
	}

	if _, _, err := c.client.Repositories.CreateHook(ctx, c.owner, c.repo, hook); err != nil {
		// Ignore exists error if the list doesn't return all installed hooks
		if strings.Contains(err.Error(), "Hook already exists on this repository") {
			return nil
		}
		return err
	}

	return nil
}

func (c *githubConfigurator) RepoInfo(ctx context.Context) (*models.RepoInfo, error) {
	repo, _, err := c.client.Repositories.Get(ctx, c.owner, c.repo)
	if err != nil {
		return nil, err
	}

	branch := *repo.DefaultBranch
	return &models.RepoInfo{
		ID:              int64(*repo.ID),
		Provider:        "github",
		RepoPath:        *repo.FullName,
		RepoBranch:      branch,
		AllowedBranches: []string{branch},
	}, nil
}
