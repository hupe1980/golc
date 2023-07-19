package embedding

import (
	"context"

	"cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Palm satisfies the Embedder interface.
var _ schema.Embedder = (*Palm)(nil)

// PalmClient is an interface for the Palm client.
type PalmClient interface {
	EmbedText(context.Context, *generativelanguagepb.EmbedTextRequest, ...gax.CallOption) (*generativelanguagepb.EmbedTextResponse, error)
}

// PalmOptions contains options for configuring the Palm client.
type PalmOptions struct {
	ModelName string
}

// Palm is a client for the Palm embedding service.
type Palm struct {
	client PalmClient
	opts   PalmOptions
}

// NewPalm creates a new instance of the Palm client.
func NewPalm(client PalmClient, optFns ...func(o *PalmOptions)) *Palm {
	opts := PalmOptions{
		ModelName: "models/embedding-gecko-001",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Palm{
		client: client,
		opts:   opts,
	}
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *Palm) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))

	for i, text := range texts {
		v, err := e.EmbedQuery(ctx, text)
		if err != nil {
			return nil, err
		}

		embeddings[i] = v
	}

	return embeddings, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *Palm) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	res, err := e.client.EmbedText(ctx, &generativelanguagepb.EmbedTextRequest{
		Model: e.opts.ModelName,
		Text:  text,
	})
	if err != nil {
		return nil, err
	}

	values := res.GetEmbedding().GetValue()

	embedding := make([]float64, len(values))
	for i, v := range values {
		embedding[i] = float64(v)
	}

	return embedding, nil
}
