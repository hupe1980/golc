package schema

import "context"

type Document struct {
	PageContent string
	Metadata    map[string]any
}

type DocumentLoader interface {
	Load(context.Context) ([]Document, error)
	LoadAndSplit(ctx context.Context, splitter TextSplitter)
}

type Retriever interface {
	GetRelevantDocuments(ctx context.Context, query string) ([]Document, error)
}

type TextSplitter interface {
	SplitDocuments(docs []Document) ([]Document, error)
}
