package embedding

import (
	"context"

	"github.com/hupe1980/golc/integration/ernie"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Ernie satisfies the Embedder interface.
var _ schema.Embedder = (*Ernie)(nil)

// ErnieClient is an interface for interacting with the Ernie API for text embedding.
type ErnieClient interface {
	// CreateEmbedding generates text embeddings using the specified model and request.
	CreateEmbedding(ctx context.Context, model string, request ernie.EmbeddingRequest) (*ernie.EmbeddingResponse, error)
}

// ErnieOptions represents configuration options for the Ernie text embedding component.
type ErnieOptions struct {
	Model string
}

// Ernie represents the text embedding component powered by Ernie.
type Ernie struct {
	client ErnieClient
	opts   ErnieOptions
}

// NewErnie creates a new instance of the Ernie text embedding component with default options.
func NewErnie(clientID, clientSecret string, optFns ...func(o *ErnieOptions)) *Ernie {
	client := ernie.New(clientID, clientSecret)

	return NewErnieFromClient(client, optFns...)
}

// NewErnieFromClient creates a new instance of the Ernie text embedding component with a custom ErnieClient and optional configuration.
func NewErnieFromClient(client ErnieClient, optFns ...func(o *ErnieOptions)) *Ernie {
	opts := ErnieOptions{
		Model: "ernie-text-embedding",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Ernie{
		client: client,
		opts:   opts,
	}
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *Ernie) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	res, err := e.client.CreateEmbedding(ctx, e.opts.Model, ernie.EmbeddingRequest{
		Input: texts,
	})
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float64, len(res.Data))
	for i, d := range res.Data {
		embeddings[i] = d.Embedding
	}

	return embeddings, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *Ernie) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	res, err := e.client.CreateEmbedding(ctx, e.opts.Model, ernie.EmbeddingRequest{
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}

	return res.Data[0].Embedding, nil
}
