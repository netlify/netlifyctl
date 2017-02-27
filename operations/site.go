package operations

import (
	"fmt"

	"strconv"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/netlify/open-api/go/porcelain/context"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var AssumeYes bool

func ConfirmCreateSite(cmd *cobra.Command) bool {
	if AssumeYes {
		return true
	}

	return askForConfirmation("We cannot find a site for this repository, do you want to create a new one?")
}

func CreateSite(cmd *cobra.Command, client *porcelain.Netlify, ctx context.Context) (*models.Site, error) {
	domain, err := AskForInput("Type your domain or press enter to use a Netlify subdomain:", "", validateCustomDomain)
	if err != nil {
		return nil, err
	}

	newS := &models.Site{
		CustomDomain: domain,
	}

	// Only configure DNS and TLS for custom domains.
	// Netlify hosted sites don't need DNS entries
	// and the connection is always over TLS.
	withTLS := len(domain) > 0

	site, err := client.CreateSite(ctx, newS, withTLS)
	if err != nil {
		return nil, err
	}

	if withTLS {
		fmt.Println("=> Provisioning TLS certificate with Let's Encrypt")

		cert, err := client.ConfigureSiteTLSCertificate(ctx, site.ID, nil)
		if err != nil {
			return site, err
		}

		cert, err = client.WaitUntilTLSCertificateReady(ctx, site.ID, cert)
		if err != nil {
			return site, err
		}

		site.ForceSsl = true
		_, err = client.UpdateSite(ctx, site)
		if err != nil {
			return site, err
		}
	}

	return site, nil
}

func ChooseSite(client *porcelain.Netlify, ctx context.Context) (*models.Site, error) {
	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		return nil, err
	}
	nameToId := make(map[string]int)
	for i, s := range sites {
		fmt.Printf("[%d] %s\n", i+1, s.Name)
		nameToId[s.Name] = i
	}

	var id int
	AskForInput("Which site?", "", func(input string) error {
		var err error

		if index, ok := nameToId[input]; ok {
			id = index
			return nil
		}

		id, err = strconv.Atoi(input)
		if err != nil {
			return err
		}

		if id > len(sites) || id <= 0 {
			return errors.New("Not a valid site selection")
		}

		return nil
	})

	return sites[id], nil
}
