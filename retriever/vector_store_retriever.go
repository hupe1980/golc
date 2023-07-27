package retriever

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure VectorStore satisfies the Retriever interface.
var _ schema.Retriever = (*VectorStore)(nil)

type VectorStoreSearchType string

const (
	VectorStoreSearchTypeSimilarity VectorStoreSearchType = "similarity"
)

type VectorStoreOptions struct {
	*schema.CallbackOptions
	SearchType VectorStoreSearchType
}

type VectorStore struct {
	v    schema.VectorStore
	opts VectorStoreOptions
}

func NewVectorStore(vectorStore schema.VectorStore, optFns ...func(o *VectorStoreOptions)) *VectorStore {
	opts := VectorStoreOptions{
		SearchType: VectorStoreSearchTypeSimilarity,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
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

// Verbose returns the verbosity setting of the retriever.
func (r *VectorStore) Verbose() bool {
	return r.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the retriever.
func (r *VectorStore) Callbacks() []schema.Callback {
	return r.opts.CallbackOptions.Callbacks
}
