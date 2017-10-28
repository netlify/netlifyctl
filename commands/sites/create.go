package sites

import (
	"errors"
	"fmt"

	"github.com/netlify/netlifyctl/commands/middleware"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/operations"
	"github.com/netlify/open-api/go/models"
	"github.com/spf13/cobra"
)

type siteCreateCmd struct {
	name         string
	customDomain string
	password     string
	forceTLS     bool
	sessionID    string
}

func setupCreateCommand(middlewares []middleware.Middleware) *cobra.Command {
	cmd := &siteCreateCmd{}
	ccmd := &cobra.Command{
		Use:   "create <-n NAME> ...",
		Short: "create site",
		Long:  "create site",
	}
	ccmd.Flags().StringVarP(&cmd.name, "name", "n", "", "site's Netlify name/subdomain")
	ccmd.Flags().StringVarP(&cmd.customDomain, "custom-domain", "c", "", "site's custom domain")
	ccmd.Flags().StringVarP(&cmd.password, "password", "p", "", "site's access password")
	ccmd.Flags().BoolVarP(&cmd.forceTLS, "force-tls", "t", false, "force TLS connections")
	ccmd.Flags().StringVarP(&cmd.sessionID, "session-id", "s", "", "Session ID for later site transfers")

	return middleware.SetupCommand(ccmd, cmd.createSite, middlewares)
}

func (c *siteCreateCmd) createSite(ctx context.Context, cmd *cobra.Command, args []string) error {
	var configFile = cmd.Root().Flag("config").Value.String()
	var conf, err = configuration.Load(configFile)
	if err != nil {
		return err
	}
	if conf.Settings.ID != "" {
		if !operations.ConfirmOverwriteSite(cmd) {
			return errors.New("Canceled")
		}
	}

	site := &models.Site{
		CustomDomain: c.customDomain,
		Name:         c.name,
		Password:     c.password,
		Ssl:          c.forceTLS,
		SessionID:    c.sessionID,
	}
	client := context.GetClient(ctx)

	fmt.Println("Creating site")
	site, err = operations.CreateSite(client, ctx, site)

	if err == nil {
		conf.Settings.ID = site.ID
		configuration.Save(configFile, conf)
		fmt.Printf("=> Done, create website with %s\n", site.ID)
	}

	return err
}
