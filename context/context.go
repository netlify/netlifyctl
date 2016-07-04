package context

import (
	"github.com/docker/distribution/context"
	"github.com/netlify/open-api/go/porcelain"
)

type Context interface {
	context.Context
}

func WithClient(ctx context.Context, client *porcelain.Netlify) context.Context {
	return context.WithValue(ctx, "netlify_client", client)
}

func GetClient(ctx context.Context) *porcelain.Netlify {
	return ctx.Value("netlify_client").(*porcelain.Netlify)
}
