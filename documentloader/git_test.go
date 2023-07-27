package documentloader

import (
	"context"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	t.Run("TestCloneURL", func(t *testing.T) {
		git, err := NewGitFromCloneURL("https://github.com/hupe1980/golc")
		require.NoError(t, err)

		docs, err := git.Load(context.Background())
		require.NoError(t, err)
		require.Greater(t, len(docs), 1)
	})

	t.Run("TestCloneURLWithFilter", func(t *testing.T) {
		git, err := NewGitFromCloneURL("https://github.com/hupe1980/golc", func(o *GitCloneURLOptions) {
			o.FileFilter = func(f *object.File) bool { return false }
		})
		require.NoError(t, err)

		docs, err := git.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 0, len(docs))
	})
}
