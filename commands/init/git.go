package init

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
	"gopkg.in/src-d/go-git.v4"
)

const (
	githubHost = "github.com"
	gitlabHost = "gitlab.com"
)

var gitSSHURL = regexp.MustCompile("\\w+@([^:]+):([^.]+)(\\.git)?")

type configurator interface {
	SetupDeployKey(context.Context, *models.DeployKey) error
	SetupWebHook(context.Context, *models.Site) error
	RepoInfo(ctx context.Context) (*models.RepoInfo, error)
}

type gitRepo struct {
	Remote        string
	URL           *url.URL
	CurrentBranch string
}

func loadConfigurator(ctx context.Context, provider *url.URL) (configurator, error) {
	switch {
	case provider.Host == githubHost:
		return newGitHubConfigurator(ctx, provider)
	case provider.Host == gitlabHost:
		return newGitLabConfigurator(ctx, provider)
	}
	return nil, fmt.Errorf("git provider not supported: %s", provider)
}

func getRepo() (*gitRepo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return getRepoFromPath(cwd)
}

func getRepoFromPath(cwd string) (*gitRepo, error) {
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, err
	}

	u := remote.Config().URLs[0]
	p, err := url.Parse(u)
	if err != nil {
		if matches := gitSSHURL.FindStringSubmatch(u); len(matches) > 0 {
			p = &url.URL{
				Host: matches[1],
				Path: matches[2],
			}
		} else {
			return nil, err
		}
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	branch := strings.Split(string(head.Name()), "/")
	return &gitRepo{
		Remote:        u,
		URL:           p,
		CurrentBranch: branch[len(branch)-1],
	}, nil
}
