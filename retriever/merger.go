package retriever

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Merger satisfies the Retriever interface.
var _ schema.Retriever = (*Merger)(nil)

type MergerOptions struct {
	*schema.CallbackOptions
}

type Merger struct {
	retrievers []schema.Retriever
	opts       MergerOptions
}

func NewMerger(retrievers []schema.Retriever, optFns ...func(o *MergerOptions)) *Merger {
	opts := MergerOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

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

// Verbose returns the verbosity setting of the retriever.
func (r *Merger) Verbose() bool {
	return r.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the retriever.
func (r *Merger) Callbacks() []schema.Callback {
	return r.opts.CallbackOptions.Callbacks
}
