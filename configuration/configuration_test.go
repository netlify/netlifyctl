package configuration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	c, err := Load("test/netlify-headers-test.toml")
	require.NoError(t, err)
	require.Equal(t, "siteid", c.Settings.ID)
	require.Len(t, c.Redirects, 1)
	require.Equal(t, "BAR", c.Redirects[0].Headers["FOO"])
}
