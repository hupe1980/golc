package documentcompressor

import (
	"context"

	cohere "github.com/cohere-ai/cohere-go/v2"
	core "github.com/cohere-ai/cohere-go/v2/core"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure CohereRerank satisfies the DocumentCompressor interface.
var _ schema.DocumentCompressor = (*CohereRerank)(nil)

// CohereClient is an interface for interacting with the Cohere API.
type CohereClient interface {
	Rerank(ctx context.Context, request *cohere.RerankRequest, opts ...core.RequestOption) (*cohere.RerankResponse, error)
}

// CohereRerankOptions contains options for Cohere Rerank compression.
type CohereRerankOptions struct {
	ModelName string
	TopN      int
}

// CohereRerank is a struct representing the Cohere Rerank compression functionality.
type CohereRerank struct {
	client CohereClient
	opts   CohereRerankOptions
}

// NewCohereRank creates a new instance of CohereRerank with the provided client and options.
func NewCohereRank(client CohereClient, optFns ...func(o *CohereRerankOptions)) *CohereRerank {
	opts := CohereRerankOptions{
		ModelName: "rerank-multilingual-v2.0",
		TopN:      3,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &CohereRerank{
		client: client,
		opts:   opts,
	}
}

// Compress compresses the input documents using Cohere Rerank.
func (c *CohereRerank) Compress(ctx context.Context, docs []schema.Document, query string) ([]schema.Document, error) {
	items := make([]*cohere.RerankRequestDocumentsItem, len(docs))

	for i, doc := range docs {
		items[i] = &cohere.RerankRequestDocumentsItem{
			String: doc.PageContent,
		}
	}

	res, err := c.client.Rerank(ctx, &cohere.RerankRequest{
		Model:           util.AddrOrNil(c.opts.ModelName),
		Documents:       items,
		Query:           query,
		TopN:            util.AddrOrNil(c.opts.TopN),
		ReturnDocuments: util.AddrOrNil(false),
	})
	if err != nil {
		return nil, err
	}

	compressedDocs := make([]schema.Document, len(res.Results))

	for i, r := range res.Results {
		compressedDocs[i] = docs[r.Index]

		if len(compressedDocs[i].Metadata) > 0 {
			compressedDocs[i].Metadata["relevanceScore"] = r.RelevanceScore
		} else {
			compressedDocs[i].Metadata = map[string]any{
				"relevanceScore": r.RelevanceScore,
			}
		}
	}

	return compressedDocs, nil
}
