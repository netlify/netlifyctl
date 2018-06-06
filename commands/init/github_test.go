package init

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGitHubConfiguratorWithHTTPSURL(t *testing.T) {
	u, err := url.Parse("https://github.com/foo/bar")
	require.NoError(t, err)

	g, err := newGitHubConfiguratorWithClient(context.Background(), u, nil)
	require.NoError(t, err)
	require.Equal(t, "foo", g.owner)
	require.Equal(t, "bar", g.repo)
}
