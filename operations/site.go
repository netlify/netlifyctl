package operations

import (
	"fmt"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/netlify/open-api/go/porcelain/context"
	"github.com/spf13/cobra"
)

var AssumeYes bool

func ConfirmCreateSite(cmd *cobra.Command) bool {
	if AssumeYes {
		fmt.Println("Creating new site")
		return true
	}

	return askForConfirmation("We cannot find a site for this repository, do you want to create a new one?")
}

func CreateSite(cmd *cobra.Command, client *porcelain.Netlify, ctx context.Context) (*models.Site, error) {
	domain, err := askForInput("Type your domain or press enter to use a Netlify subdomain: ", "", validateCustomDomain)
	if err != nil {
		return nil, err
	}

	newS := &models.Site{
		CustomDomain: domain,
	}

	site, err := client.CreateSite(ctx, newS)
	if err != nil {
		return nil, err
	}

	if len(domain) != 0 {
		// TODO: Register DNS with a DNS provider and Configure Let's Encrypt
	}

	return site, nil
}
