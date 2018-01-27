package sites

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/open-api/go/models"
	"github.com/spf13/cobra"
)

const (
	defaultEditor = "vi"
	message       = `
# Please modify only the fields you want to change.
# Lines starting with '#' will be ignored, and removing all content will
# cancel the update.
#`
)

type siteUpdateCmd struct {
	name         string
	customDomain string
	password     string
	forceTLS     bool
}

type editableSite struct {
	Name          string
	CustomDomain  string
	Password      string
	ForceTLS      bool
	DomainAliases []string
}

func setupUpdateCommand(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &siteUpdateCmd{}
	ccmd := &cobra.Command{
		Use:   "update <-n NAME> ...",
		Short: "Update site settings",
		Long:  "Update site settings",
	}
	ccmd.Flags().StringVarP(&cmd.name, "name", "n", "", "site's Netlify name/subdomain")
	ccmd.Flags().StringVarP(&cmd.customDomain, "custom-domain", "c", "", "site's custom domain")
	ccmd.Flags().StringVarP(&cmd.password, "password", "p", "", "site's access password")
	ccmd.Flags().BoolVarP(&cmd.forceTLS, "force-tls", "t", false, "force TLS connections")

	siteMiddlewares := append([]middleware.Middleware{middleware.SiteConfigMiddleware}, middlewares...)

	return middleware.SetupCommand(ccmd, cmd.updateSite, siteMiddlewares)
}

func (c *siteUpdateCmd) updateSite(ctx context.Context, cmd *cobra.Command, args []string) error {
	conf := context.GetSiteConfig(ctx)
	if conf.Settings.ID == "" {
		return errors.New("Failed to load site configuration")
	}

	client := context.GetClient(ctx)
	site, err := client.GetSite(ctx, conf.Settings.ID)
	if err != nil {
		return err
	}

	if c.showEditor(site) {
		edit, err := openEditor(site)
		if err != nil {
			return err
		}
		site.Name = edit.Name
		site.CustomDomain = edit.CustomDomain
		site.Password = edit.Password
		site.Ssl = edit.ForceTLS
		site.DomainAliases = edit.DomainAliases
	} else {
		if c.name != "" {
			site.Name = c.name
		}
		if c.customDomain != "" {
			site.CustomDomain = c.customDomain
		}
		if c.password != "" {
			site.Password = c.password
		}
		site.Ssl = c.forceTLS
	}

	if site.Ssl {
		site.ForceSsl = true
	}

	setup := &models.SiteSetup{Site: *site}
	updated, err := client.UpdateSite(ctx, setup)
	if err != nil {
		return err
	}
	fmt.Printf("site updated: %s", updated.URL)

	return nil
}

func (c *siteUpdateCmd) showEditor(site *models.Site) bool {
	return c.name == "" && c.customDomain == "" && c.password == "" && c.forceTLS == site.Ssl
}

func openEditor(site *models.Site) (*editableSite, error) {
	tmpDir := os.TempDir()
	tmpFile, err := ioutil.TempFile(tmpDir, "netlifyctl-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	es := &editableSite{
		Name:          site.Name,
		CustomDomain:  site.CustomDomain,
		ForceTLS:      site.Ssl,
		Password:      site.Password,
		DomainAliases: site.DomainAliases,
	}

	if err := toml.NewEncoder(tmpFile).Encode(es); err != nil {
		return nil, err
	}

	if _, err := tmpFile.WriteString(message); err != nil {
		return nil, err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		// Cancel editing if content is empty
		return nil, nil
	}

	if _, err := toml.Decode(string(b), es); err != nil {
		return nil, err
	}

	return es, nil
}
