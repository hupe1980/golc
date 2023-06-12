package documentloader

import (
	"context"
	"io"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/textsplitter"
)

type TextLoader struct {
	r io.Reader
}

func NewTextLoader(r io.Reader) *TextLoader {
	return &TextLoader{
		r: r,
	}
}

func (l *TextLoader) Load(ctx context.Context) ([]golc.Document, error) {
	b, err := io.ReadAll(l.r)
	if err != nil {
		return nil, err
	}

	return []golc.Document{
		{
			PageContent: string(b),
			Metadata:    map[string]any{},
		},
	}, nil
}

func (l *TextLoader) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]golc.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
