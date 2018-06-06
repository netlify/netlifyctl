package init

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
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
	url = formatGitLabURL(url)

	username, err := getUsername(url.Host)
	if err != nil {
		return nil, err
	}

	password, err := getPassword(url.Host)
	if err != nil {
		return nil, err
	}

	u := strings.TrimSpace(username)
	p := strings.TrimSpace(password)
	if u == "" || p == "" {
		return nil, errors.New("Username and password cannot be blank")
	}

	client, err := gitlab.NewBasicAuthClient(nil, url.Scheme+"://"+url.Host, u, p)
	if err != nil {
		return nil, err
	}

	return &gitlabConfigurator{
		client:   client,
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

func (c *gitlabConfigurator) RepoInfo(ctx context.Context) (*models.RepoInfo, error) {
	project, _, err := c.client.Projects.GetProject(c.repoPath)
	if err != nil {
		return nil, err
	}

	branch := project.DefaultBranch
	return &models.RepoInfo{
		ID:              int64(project.ID),
		Provider:        "gitlab",
		RepoPath:        project.PathWithNamespace,
		RepoBranch:      branch,
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
func getPassword(host string) (string, error) {
	fmt.Printf("%s password: ", host)

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

func getSession(urlStr string, opt *gitlab.GetSessionOptions, options ...gitlab.OptionFunc) (*gitlab.Session, *gitlab.Response, error) {
	client := gitlab.NewClient(nil, "")
	client.SetBaseURL(urlStr)
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

func formatGitLabURL(u *url.URL) *url.URL {
	u.Scheme = "https"
	u.Path = strings.TrimRight(strings.TrimLeft(u.Path, "/"), ".git")

	return u
}
