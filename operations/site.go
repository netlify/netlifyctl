package operations

import (
	"fmt"

	"strconv"

	"os"

	"errors"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/netlify/open-api/go/porcelain/context"
	"github.com/spf13/cobra"
)

var AssumeYes bool

func ConfirmCreateSite(cmd *cobra.Command) bool {
	if AssumeYes {
		return true
	}

	return askForConfirmation("We cannot find a site for this repository, do you want to create a new one?")
}

func CreateSite(client *porcelain.Netlify, ctx context.Context) (*models.Site, error) {
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

func ChooseOrCreateSite(client *porcelain.Netlify, ctx context.Context) (*models.Site, error) {
	fmt.Println("No site configured in the netlify.toml, fetching your existing sites.")
	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Choose a site to deploy to or 0 to create a new site.")
	nameToId := make(map[string]*models.Site)
	fmt.Println("[0] Create a new site")
	for i, s := range sites {
		fmt.Printf("[%d] %s\n", i+1, s.Name)
		nameToId[s.Name] = s
	}

	for {
		input, err := AskForInput("Which site?", "0", nil)
		if err == nil {
			if selection, ok := nameToId[input]; ok {
				return selection, nil
			}

			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Fprint(os.Stdout, "Input must be an index or the site name")
				continue
			}

			if id == 0 {
				// in this case we want to do whatever the create says
				// that includes storing off the new id
				return CreateSite(client, ctx)
			}

			if id > len(sites) || id < 0 {
				fmt.Fprint(os.Stdout, "Input must be an index or the site name")
			}

			return sites[id], nil
		}
	}

	return nil, errors.New("Failed to select a site")
}
