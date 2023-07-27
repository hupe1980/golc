// Package retriever provides functionality for retrieving relevant documents using various services.
package retriever

import (
	"context"
	"net/http"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	Callbacks   []schema.Callback
	ParentRunID string
}

func Run(ctx context.Context, retriever schema.Retriever, query string, optFns ...func(*Options)) ([]schema.Document, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, retriever.Callbacks(), retriever.Verbose(), func(mo *callback.ManagerOptions) {
		mo.ParentRunID = opts.ParentRunID
	})

	rm, err := cm.OnRetrieverStart(ctx, &schema.RetrieverStartManagerInput{
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	docs, err := retriever.GetRelevantDocuments(ctx, query)
	if err != nil {
		if cbErr := rm.OnRetrieverError(ctx, &schema.RetrieverErrorManagerInput{
			Error: err,
		}); cbErr != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if err := rm.OnRetrieverEnd(ctx, &schema.RetrieverEndManagerInput{
		Docs: docs,
	}); err != nil {
		return nil, err
	}

	return docs, nil
}
