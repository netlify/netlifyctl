package init

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
	gitlab "github.com/xanzy/go-gitlab"
	"golang.org/x/crypto/ssh/terminal"
)

type gitlabConfigurator struct {
	client   *gitlab.Client
	repoPath string
}

func newGitLabConfigurator(ctx context.Context, url *url.URL) (*gitlabConfigurator, error) {
	username, err := getUsername(url.Host)
	if err != nil {
		return nil, err
	}

	password, err := getPassword()
	if err != nil {
		return nil, err
	}

	session, _, err := getSession(&gitlab.GetSessionOptions{
		Login:    gitlab.String(username),
		Password: gitlab.String(password),
	}, nil)
	if err != nil {
		return nil, err
	}

	return &gitlabConfigurator{
		client:   gitlab.NewClient(nil, session.PrivateToken),
		repoPath: url.Path,
	}, nil
}

func (c *gitlabConfigurator) SetupDeployKey(ctx context.Context, deployKey *models.DeployKey) error {
	key := &gitlab.AddDeployKeyOptions{
		Title: gitlab.String(deployKeyTitle),
		Key:   gitlab.String(deployKey.PublicKey),
	}
	_, _, err := c.client.DeployKeys.AddDeployKey(c.repoPath, key)
	return err
}

func (c *gitlabConfigurator) SetupWebHook(ctx context.Context, site *models.Site) error {
	hook := &gitlab.AddProjectHookOptions{
		URL:                 gitlab.String(site.DeployHook),
		PushEvents:          gitlab.Bool(true),
		MergeRequestsEvents: gitlab.Bool(true),
	}
	_, _, err := c.client.Projects.AddProjectHook(c.repoPath, hook)
	return err
}

func (c *gitlabConfigurator) RepoInfo(ctx context.Context) (*models.RepoSetup, error) {
	project, _, err := c.client.Projects.GetProject(c.repoPath)
	if err != nil {
		return nil, err
	}

	branch := project.DefaultBranch
	return &models.RepoSetup{
		ID:              int64(project.ID),
		Provider:        "gitlab",
		Repo:            project.Path,
		Branch:          branch,
		AllowedBranches: []string{branch},
	}, nil
}

func getUsername(host string) (string, error) {
	fmt.Printf("%s username: ", host)

	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return line, nil
}

// copied from GitHub's hub
func getPassword() (string, error) {
	stdin := int(syscall.Stdin)
	initialTermState, err := terminal.GetState(stdin)
	if err != nil {
		return "", err
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		terminal.Restore(stdin, initialTermState)
		switch sig := s.(type) {
		case syscall.Signal:
			if int(sig) == 2 {
				fmt.Println("^C")
			}
		}
		os.Exit(1)
	}()

	passBytes, err := terminal.ReadPassword(stdin)
	if err != nil {
		return "", err
	}

	signal.Stop(c)
	fmt.Print("\n")
	return string(passBytes), nil
}

func getSession(opt *gitlab.GetSessionOptions, options ...gitlab.OptionFunc) (*gitlab.Session, *gitlab.Response, error) {
	client := &gitlab.Client{}
	req, err := client.NewRequest("POST", "session", opt, options)
	if err != nil {
		return nil, nil, err
	}

	session := new(gitlab.Session)
	resp, err := client.Do(req, session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, err
}
