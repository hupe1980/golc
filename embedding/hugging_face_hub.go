package embedding

import (
	"context"

	huggingface "github.com/hupe1980/go-huggingface"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure HuggingFaceHub satisfies the Embedder interface.
var _ schema.Embedder = (*HuggingFaceHub)(nil)

// HuggingFaceHubClient represents a client for interacting with Hugging Face Hub.
type HuggingFaceHubClient interface {
	// FeatureExtractionWithAutomaticReduction performs feature extraction with automatic reduction.
	// It returns the extraction response or an error if the operation fails.
	FeatureExtractionWithAutomaticReduction(ctx context.Context, req *huggingface.FeatureExtractionRequest) (huggingface.FeatureExtractionWithAutomaticReductionResponse, error)
}

type HuggingFaceHubOptions struct {
	// Model to use for embedding.
	Model string
	// Options represents optional settings for the feature extraction.
	Options huggingface.Options
}

// HuggingFaceHub represents an embedder for Hugging Face Hub models.
type HuggingFaceHub struct {
	client HuggingFaceHubClient
	opts   HuggingFaceHubOptions
}

// NewHuggingFaceHub creates a new instance of the HuggingFaceHub embedder.
func NewHuggingFaceHub(token string, optFns ...func(o *HuggingFaceHubOptions)) *HuggingFaceHub {
	client := huggingface.NewInferenceClient(token)

	return NewHuggingFaceHubFromClient(client, optFns...)
}

// NewHuggingFaceHubFromClient creates a new instance of the HuggingFaceHub embedder from a custom client.
func NewHuggingFaceHubFromClient(client HuggingFaceHubClient, optFns ...func(o *HuggingFaceHubOptions)) *HuggingFaceHub {
	opts := HuggingFaceHubOptions{
		Model:   "sentence-transformers/all-mpnet-base-v2",
		Options: huggingface.Options{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &HuggingFaceHub{
		client: client,
		opts:   opts,
	}
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *HuggingFaceHub) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	res, err := e.client.FeatureExtractionWithAutomaticReduction(ctx, &huggingface.FeatureExtractionRequest{
		Inputs:  texts,
		Model:   e.opts.Model,
		Options: e.opts.Options,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *HuggingFaceHub) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	res, err := e.client.FeatureExtractionWithAutomaticReduction(ctx, &huggingface.FeatureExtractionRequest{
		Inputs:  []string{text},
		Model:   e.opts.Model,
		Options: e.opts.Options,
	})
	if err != nil {
		return nil, err
	}

	return res[0], nil
}
