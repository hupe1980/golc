package schema

import "context"

type Document struct {
	PageContent string
	Metadata    map[string]any
}

type DocumentLoader interface {
	Load(ctx context.Context) ([]Document, error)
	LoadAndSplit(ctx context.Context, splitter TextSplitter) ([]Document, error)
}

type DocumentCompressor interface {
	// Compress compresses the input documents.
	Compress(ctx context.Context, docs []Document, query string) ([]Document, error)
}

type DocumentTransformer interface {
	Transform(ctx context.Context, docs []Document) ([]Document, error)
}

type Retriever interface {
	GetRelevantDocuments(ctx context.Context, query string) ([]Document, error)
	// Verbose returns the verbosity setting of the retriever.
	Verbose() bool
	// Callbacks returns the registered callbacks of the retriever.
	Callbacks() []Callback
}

type TextSplitter interface {
	SplitDocuments(docs []Document) ([]Document, error)
}

type VectorStore interface {
	AddDocuments(ctx context.Context, docs []Document) error
	SimilaritySearch(ctx context.Context, query string) ([]Document, error)
}
