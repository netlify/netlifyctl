package commands

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/netlify/netlifyctl/auth"
	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your site",
	Long:  "Deploy your site",
	RunE:  deploySite,
}

func deploySite(cmd *cobra.Command, args []string) error {
	site, err := configuration.Load()
	if err != nil {
		return err
	}

	// TODO: A new site won't ever have an ID.
	//       Ask to create a new site and save ID.

	c := porcelain.NewHTTPClient(nil)

	logrus.WithFields(logrus.Fields{"site": site.Settings.ID, "root": site.Settings.Root()}).Debug("deploy site")

	resp, err := c.DeploySite(site.Settings.ID, site.Settings.Root(), auth.ClientCredentials())
	if err != nil {
		return err
	}
	fmt.Println(resp.SiteURL)

	return nil
}
