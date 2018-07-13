package init

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoInfo(t *testing.T) {
	t.Run("repo info with http remote", func(t *testing.T) {
		i := &gitRepo{
			Remote:        "https://github.com/calavera/test",
			URL:           &url.URL{Host: "github.com", Path: "calavera/test"},
			CurrentBranch: "master",
		}

		m := manualConfigurator{i}
		r, err := m.RepoInfo(context.Background())
		require.NoError(t, err)

		require.Equal(t, "git@github.com:calavera/test.git", r.RepoPath)
		require.Equal(t, "manual", r.Provider)
		require.Equal(t, "master", r.RepoBranch)
		require.ElementsMatch(t, []string{"master"}, r.AllowedBranches)
	})
}

func TestRepoInfoWithGitExtension(t *testing.T) {
	t.Run("repo info with http remote", func(t *testing.T) {
		i := &gitRepo{
			Remote:        "https://github.com/calavera/test.git",
			URL:           &url.URL{Host: "github.com", Path: "calavera/test"},
			CurrentBranch: "master",
		}

		m := manualConfigurator{i}
		r, err := m.RepoInfo(context.Background())
		require.NoError(t, err)

		require.Equal(t, "git@github.com:calavera/test.git", r.RepoPath)
		require.Equal(t, "manual", r.Provider)
		require.Equal(t, "master", r.RepoBranch)
		require.ElementsMatch(t, []string{"master"}, r.AllowedBranches)
	})
}
