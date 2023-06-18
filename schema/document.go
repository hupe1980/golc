package schema

import "context"

type Document struct {
	PageContent string
	Metadata    map[string]any
}

type DocumentLoader interface {
	Load(ctx context.Context) ([]Document, error)
	LoadAndSplit(ctx context.Context, splitter TextSplitter)
}

type Retriever interface {
	GetRelevantDocuments(ctx context.Context, query string) ([]Document, error)
}

type TextSplitter interface {
	SplitDocuments(docs []Document) ([]Document, error)
}

type VectorStore interface {
	AddDocuments(ctx context.Context, docs []Document) error
	SimilaritySearch(ctx context.Context, query string) ([]Document, error)
}
