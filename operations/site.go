package operations

import (
	"errors"
	"fmt"

	"strconv"

	"github.com/netlify/netlifyctl/context"
	"github.com/netlify/netlifyctl/ui"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain"
	"github.com/spf13/cobra"
)

func CreateSite(client *porcelain.Netlify, ctx context.Context, newS *models.Site) (*models.Site, error) {
	// Only configure DNS and TLS for custom domains.
	// Netlify hosted sites don't need DNS entries
	// and the connection is always over TLS.
	withTLS := len(newS.CustomDomain) > 0

	setup := &models.SiteSetup{Site: *newS}
	var (
		site *models.Site
		err  error
	)

	ui.Track("Creating site ... ", "Site created", func() error {
		site, err = client.CreateSite(ctx, setup, withTLS)
		return err
	})

	if err != nil {
		return nil, err
	}

	if withTLS {
		err := ui.Track("Provisioning TLS certificate with Let's Encrypt ... ", "TLS Certificate provisioned", func() error {
			cert, err := client.ConfigureSiteTLSCertificate(ctx, site.ID, nil)
			if err != nil {
				return err
			}

			cert, err = client.WaitUntilTLSCertificateReady(ctx, site.ID, cert)
			if err != nil {
				return err
			}

			site.ForceSsl = true
			setup := &models.SiteSetup{Site: *site}
			_, err = client.UpdateSite(ctx, setup)
			return err
		})
		if err != nil {
			return site, err
		}
	}

	return site, nil
}

func ChooseOrCreateSite(ctx context.Context, cmd *cobra.Command) (*models.Site, error) {
	var name string
	if f := cmd.Flag("name"); f != nil {
		name = f.Value.String()
	}

	client := context.GetClient(ctx)
	if ui.AskForConfirmation("Create a new site?") {
		site := &models.Site{
			Name: name,
		}
		return CreateSite(client, ctx, site)
	}

	if name == "" {
		i, err := ui.AskForInput("Search in your sites:", "")
		if err != nil {
			return nil, err
		}
		name = i
	}

	filter := "all"
	params := operations.NewListSitesParams().WithFilter(&filter).WithName(&name)
	sites, err := client.ListSites(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(sites) == 0 {
		fmt.Printf("There are no sites to deploy, create a new site before deploying or search again  %s", ui.ErrorCheck())
		return nil, errors.New("there are no sites to deploy")
	}

	if len(sites) == 1 {
		return sites[0], nil
	}

	nameToId := make(map[string]*models.Site)
	for i, s := range sites {
		fmt.Printf("[%d] %s\n", i+1, s.Name)
		nameToId[s.Name] = s
	}

	for {
		input, err := ui.AskForInput("Choose site by name or index:", "")
		if err == nil {
			if selection, ok := nameToId[input]; ok {
				return selection, nil
			}

			id, err := strconv.Atoi(input)
			if err != nil {
				continue
			}

			if id > len(sites) || id < 0 {
				continue
			}

			return sites[id-1], nil
		}
	}
}
