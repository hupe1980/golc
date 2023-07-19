package documentloader

import (
	"context"
	"io"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Text satisfies the DocumentLoader interface.
var _ schema.DocumentLoader = (*Text)(nil)

type Text struct {
	r io.Reader
}

// NewText creates a new Text document loader with the given reader.
func NewText(r io.Reader) *Text {
	return &Text{
		r: r,
	}
}

// Load reads the content from the reader and returns it as a single document.
func (l *Text) Load(ctx context.Context) ([]schema.Document, error) {
	b, err := io.ReadAll(l.r)
	if err != nil {
		return nil, err
	}

	return []schema.Document{
		{
			PageContent: string(b),
			Metadata:    map[string]any{},
		},
	}, nil
}

// LoadAndSplit reads the content from the reader and splits it into multiple documents using the provided splitter.
func (l *Text) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
