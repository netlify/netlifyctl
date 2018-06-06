package init

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatGitLabURL(t *testing.T) {
	u := &url.URL{
		Host: "gitlab.com",
		Path: "/foo/bar.git",
	}

	e := formatGitLabURL(u).String()
	require.Equal(t, "https://gitlab.com/foo/bar", e)
}
