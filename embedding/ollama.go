package embedding

import (
	"context"

	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/schema"
	"golang.org/x/sync/errgroup"
)

// Compile time check to ensure Ollama satisfies the Embedder interface.
var _ schema.Embedder = (*Ollama)(nil)

// OllamaClient is an interface for interacting with the Ollama model's embedding functionality.
type OllamaClient interface {
	CreateEmbedding(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error)
}

// OllamaOptions contains options for configuring the Ollama model.
type OllamaOptions struct {
	MaxConcurrency int
	// ModelName is the name of the Gemini model to use.
	ModelName string `map:"model_name,omitempty"`
}

// Ollama is a struct representing the Ollama embedding model.
type Ollama struct {
	client OllamaClient
	opts   OllamaOptions
}

// NewOllama creates a new instance of the Ollama embedding model.
func NewOllama(client OllamaClient, optFns ...func(o *OllamaOptions)) *Ollama {
	opts := OllamaOptions{
		MaxConcurrency: 5,
		ModelName:      "llama2",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Ollama{
		client: client,
		opts:   opts,
	}
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Ollama) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	errs, errctx := errgroup.WithContext(ctx)

	errs.SetLimit(e.opts.MaxConcurrency)

	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		i, text := i, text

		errs.Go(func() error {
			res, err := e.client.CreateEmbedding(errctx, &ollama.EmbeddingRequest{
				Prompt: text,
				Model:  e.opts.ModelName,
			})
			if err != nil {
				return err
			}

			embeddings[i] = res.Embedding

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return embeddings, nil
}

// EmbedText embeds a single text and returns its embedding.
func (e *Ollama) EmbedText(ctx context.Context, text string) ([]float32, error) {
	res, err := e.client.CreateEmbedding(ctx, &ollama.EmbeddingRequest{
		Prompt: text,
		Model:  e.opts.ModelName,
	})
	if err != nil {
		return nil, err
	}

	return res.Embedding, nil
}
