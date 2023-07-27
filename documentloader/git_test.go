package documentloader

import (
	"context"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	t.Run("TestLoad", func(t *testing.T) {
		git, err := NewGitFromPath("testdata/gitrepo")
		require.NoError(t, err)

		docs, err := git.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, "This is file a.txt.\n\n", docs[0].PageContent)
		require.Equal(t, "a.txt", docs[0].Metadata["name"].(string))
		require.Equal(t, "This is file b.txt.\n\n", docs[1].PageContent)
		require.Equal(t, "b.txt", docs[1].Metadata["name"].(string))
	})

	t.Run("TestLoadWithFilter", func(t *testing.T) {
		git, err := NewGitFromPath("testdata/gitrepo", func(o *GitOptions) {
			o.FileFilter = func(f *object.File) bool {
				return f.Name != "b.txt"
			}
		})
		require.NoError(t, err)

		docs, err := git.Load(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, len(docs))
		require.Equal(t, "This is file a.txt.\n\n", docs[0].PageContent)
		require.Equal(t, "a.txt", docs[0].Metadata["name"].(string))
	})

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
