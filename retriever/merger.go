package retriever

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Merger satisfies the Retriever interface.
var _ schema.Retriever = (*Merger)(nil)

type Merger struct {
	retrievers []schema.Retriever
}

func NewMerger(retrievers ...schema.Retriever) *Merger {
	return &Merger{
		retrievers: retrievers,
	}
}

func (r *Merger) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	var err error

	// Get the results of all retrievers.
	retrieverDocs := make([][]schema.Document, len(r.retrievers))

	for i, retriever := range r.retrievers {
		retrieverDocs[i], err = retriever.GetRelevantDocuments(ctx, query)
		if err != nil {
			return nil, err
		}
	}

	// Merge the results of the retrievers.
	mergedDocuments := make([]schema.Document, 0)

	maxDocs := 0

	for _, docs := range retrieverDocs {
		if len(docs) > maxDocs {
			maxDocs = len(docs)
		}
	}

	for i := 0; i < maxDocs; i++ {
		for j := range r.retrievers {
			if i < len(retrieverDocs[j]) {
				mergedDocuments = append(mergedDocuments, retrieverDocs[j][i])
			}
		}
	}

	return mergedDocuments, nil
}
