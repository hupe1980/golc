package documentloader

import (
	"context"
	"io"

	"github.com/hupe1980/golc/schema"
)

type TextLoader struct {
	r io.Reader
}

func NewTextLoader(r io.Reader) *TextLoader {
	return &TextLoader{
		r: r,
	}
}

func (l *TextLoader) Load(ctx context.Context) ([]schema.Document, error) {
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

func (l *TextLoader) LoadAndSplit(ctx context.Context, splitter schema.TextSplitter) ([]schema.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
