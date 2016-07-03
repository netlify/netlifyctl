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
	conf, err := configuration.Load()
	if err != nil {
		return err
	}
	s := conf.Settings

	// TODO: A new site won't ever have an ID.
	//       Ask to create a new site and save ID.

	c := porcelain.NewHTTPClient(nil)
	ctx := auth.NewContext()

	logrus.WithFields(logrus.Fields{"site": s.ID, "root": s.Root()}).Debug("deploy site")

	d, err := c.DeploySite(ctx, s.ID, s.Root())
	if err != nil {
		return err
	}

	ready, err := c.WaitUntilDeployReady(ctx, d)
	if err != nil {
		return err
	}
	fmt.Println(ready.URL)

	return nil
}
