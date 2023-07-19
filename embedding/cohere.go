package embedding

import (
	"context"

	"github.com/cohere-ai/cohere-go"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Cohere satisfies the Embedder interface.
var _ schema.Embedder = (*Cohere)(nil)

// CohereClient is an interface for the Cohere client.
type CohereClient interface {
	Embed(opts cohere.EmbedOptions) (*cohere.EmbedResponse, error)
}

// CohereOptions contains options for configuring the Cohere instance.
type CohereOptions struct {
	// Model name to use.
	Model string
	// Truncate embeddings that are too long from start or end ("NONE"|"START"|"END")
	Truncate string
}

// Cohere is a client for the Cohere API.
type Cohere struct {
	client CohereClient
	opts   CohereOptions
}

// NewCohere creates a new Cohere instance with the provided API key and options.
// It returns the initialized Cohere instance or an error if initialization fails.
func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	return NewCohereFromClient(client, optFns...)
}

// NewCohereFromClient creates a new Cohere instance from an existing Cohere client and options.
// It returns the initialized Cohere instance.
func NewCohereFromClient(client CohereClient, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		Model: "embed-english-v2.0",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Cohere{
		client: client,
	}, nil
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *Cohere) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	res, err := e.client.Embed(cohere.EmbedOptions{
		Model:    e.opts.Model,
		Truncate: e.opts.Truncate,
		Texts:    texts,
	})
	if err != nil {
		return nil, err
	}

	return res.Embeddings, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *Cohere) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	embeddings, err := e.EmbedDocuments(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return embeddings[0], nil
}
