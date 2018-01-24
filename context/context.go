package context

import (
	"context"

	"github.com/netlify/netlifyctl/configuration"
	"github.com/netlify/open-api/go/porcelain"
)

type contextKey string

type Context interface {
	context.Context
}

const (
	apiClientKey = contextKey("netlify_client")
	siteConfKey  = contextKey("netlify_site_conf")
)

func (c contextKey) String() string {
	return "netlifyctl context key " + string(c)
}

func WithClient(ctx context.Context, client *porcelain.Netlify) context.Context {
	return context.WithValue(ctx, apiClientKey, client)
}

func GetClient(ctx context.Context) *porcelain.Netlify {
	return ctx.Value(apiClientKey).(*porcelain.Netlify)
}

func WithSiteConfig(ctx context.Context, conf *configuration.Configuration) context.Context {
	return context.WithValue(ctx, siteConfKey, conf)
}

func GetSiteConfig(ctx context.Context) *configuration.Configuration {
	return ctx.Value(siteConfKey).(*configuration.Configuration)
}

func Background() context.Context {
	return context.Background()
}
