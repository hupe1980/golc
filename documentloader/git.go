package documentloader

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	memfs "github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hupe1980/golc/integration/codecommit"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Git satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*Git)(nil)

// FileFilter is a function that filters files based on specific criteria.
type FileFilter func(f *object.File) bool

// GitOptions holds options for the Git document loader.
type GitOptions struct {
	Branch     string
	FileFilter FileFilter
}

// DefaultGitOptions provides default Git options.
var DefaultGitOptions = GitOptions{
	Branch:     "main",
	FileFilter: func(f *object.File) bool { return true },
}

// GitCloneURLOptions holds options for Git repositories cloned from a URL.
type GitCloneURLOptions struct {
	GitOptions
	Auth transport.AuthMethod
}

// Git is a Git-based implementation of the DocumentLoader interface.
type Git struct {
	r    *git.Repository
	opts GitOptions
}

// NewGitFromCodeCommitURL clones a Git repository from an AWS CodeCommit URL using the provided AWS credentials,
// and returns a Git document loader. The options can be customized using functional options.
func NewGitFromCodeCommitURL(url string, creds aws.Credentials, optFns ...func(o *GitOptions)) (*Git, error) {
	signer := codecommit.NewSigner(creds)

	signedURL, err := signer.Sign(url)
	if err != nil {
		return nil, err
	}

	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: signedURL,
	})
	if err != nil {
		return nil, err
	}

	return NewGit(r, optFns...), nil
}

// NewGitFromCloneURL clones a Git repository from a URL and returns a Git document loader.
// The options can be customized using functional options.
func NewGitFromCloneURL(url string, optFns ...func(o *GitCloneURLOptions)) (*Git, error) {
	opts := GitCloneURLOptions{
		GitOptions: DefaultGitOptions,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  url,
		Auth: opts.Auth,
	})
	if err != nil {
		return nil, err
	}

	return NewGit(r, func(o *GitOptions) {
		*o = opts.GitOptions
	}), nil
}

// NewGitFromPath opens an existing Git repository from a local path and returns a Git document loader.
// The options can be customized using functional options.
func NewGitFromPath(path string, optFns ...func(o *GitOptions)) (*Git, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return NewGit(r, optFns...), nil
}

// NewGit creates a Git document loader from an existing Git repository and returns it.
// The options can be customized using functional options.
func NewGit(r *git.Repository, optFns ...func(o *GitOptions)) *Git {
	opts := DefaultGitOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Git{
		r:    r,
		opts: opts,
	}
}

// Load retrieves documents from the Git repository and returns them as a slice of schema.Document.
func (l *Git) Load(ctx context.Context) ([]schema.Document, error) {
	// Get the HEAD reference.
	ref, err := l.r.Head()
	if err != nil {
		return nil, err
	}

	// Get the commit object for the HEAD reference.
	commit, err := l.r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	docs := make([]schema.Document, 0)
	_ = tree.Files().ForEach(func(f *object.File) error {
		binary, err := f.IsBinary()
		if err != nil {
			return err
		}

		// ignore binary  or filtered files
		if binary || !l.opts.FileFilter(f) {
			return nil
		}

		contents, err := f.Contents()
		if err != nil {
			return err
		}

		docs = append(docs, schema.Document{
			PageContent: contents,
			Metadata: map[string]any{
				"name": f.Name,
			},
		})

		return nil
	})

	return docs, nil
}

// LoadAndSplit retrieves documents from the Git repository, splits them using the provided TextSplitter,
// and returns the split documents as a slice of schema.Document.
func (l *Git) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
