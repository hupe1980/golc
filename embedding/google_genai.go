package embedding

import (
	"context"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure GoogleGenAI satisfies the Embedder interface.
var _ schema.Embedder = (*GoogleGenAI)(nil)

// GoogleGenAIClient is an interface for the GoogleGenAI client.
type GoogleGenAIClient interface {
	EmbedContent(context.Context, *generativelanguagepb.EmbedContentRequest, ...gax.CallOption) (*generativelanguagepb.EmbedContentResponse, error)
	BatchEmbedContents(context.Context, *generativelanguagepb.BatchEmbedContentsRequest, ...gax.CallOption) (*generativelanguagepb.BatchEmbedContentsResponse, error)
}

// GoogleGenAIOptions contains options for configuring the GoogleGenAI client.
type GoogleGenAIOptions struct {
	ModelName string
}

// GoogleGenAI is a client for the GoogleGenAI embedding service.
type GoogleGenAI struct {
	client GoogleGenAIClient
	opts   GoogleGenAIOptions
}

// NewGoogleGenAI creates a new instance of the GoogleGenAI client.
func NewGoogleGenAI(client GoogleGenAIClient, optFns ...func(o *GoogleGenAIOptions)) *GoogleGenAI {
	opts := GoogleGenAIOptions{
		ModelName: "models/embedding-001",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &GoogleGenAI{
		client: client,
		opts:   opts,
	}
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *GoogleGenAI) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	requests := make([]*generativelanguagepb.EmbedContentRequest, len(texts))

	for i, t := range texts {
		requests[i] = &generativelanguagepb.EmbedContentRequest{
			Model: e.opts.ModelName,
			Content: &generativelanguagepb.Content{Parts: []*generativelanguagepb.Part{{
				Data: &generativelanguagepb.Part_Text{Text: t},
			}}},
		}
	}

	res, err := e.client.BatchEmbedContents(ctx, &generativelanguagepb.BatchEmbedContentsRequest{
		Model:    e.opts.ModelName,
		Requests: requests,
	})
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(texts))

	for i, e := range res.Embeddings {
		embeddings[i] = e.Values
	}

	return embeddings, nil
}

// EmbedText embeds a single text and returns its embedding.
func (e *GoogleGenAI) EmbedText(ctx context.Context, text string) ([]float32, error) {
	res, err := e.client.EmbedContent(ctx, &generativelanguagepb.EmbedContentRequest{
		Model: e.opts.ModelName,
		Content: &generativelanguagepb.Content{Parts: []*generativelanguagepb.Part{{
			Data: &generativelanguagepb.Part_Text{Text: text},
		}}},
	})
	if err != nil {
		return nil, err
	}

	return res.Embedding.Values, nil
}
