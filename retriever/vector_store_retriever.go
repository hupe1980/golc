package retriever

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure VectorStore satisfies the Retriever interface.
var _ schema.Retriever = (*VectorStore)(nil)

type VectorStoreSearchType string

const (
	VectorStoreSearchTypeSimilarity VectorStoreSearchType = "similarity"
)

type VectorStoreOptions struct {
	SearchType VectorStoreSearchType
}

type VectorStore struct {
	v    schema.VectorStore
	opts VectorStoreOptions
}

func NewVectorStore(vectorStore schema.VectorStore, optFns ...func(o *VectorStoreOptions)) *VectorStore {
	opts := VectorStoreOptions{
		SearchType: VectorStoreSearchTypeSimilarity,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &VectorStore{
		v:    vectorStore,
		opts: opts,
	}
}

// GetRelevantDocuments returns documents using the vector store.
func (r *VectorStore) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	return r.v.SimilaritySearch(ctx, query)
}
