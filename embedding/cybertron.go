package embedding

import (
	"context"

	"github.com/hupe1980/golc/schema"
	"github.com/nlpodyssey/cybertron/pkg/models/bert"
	"github.com/nlpodyssey/cybertron/pkg/tasks"
	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
)

// Compile time check to ensure Cybertron satisfies the Embedder interface.
var _ schema.Embedder = (*Cybertron)(nil)

// CybertronFromEncoderOption represents options for the Cybertron embedder.
type CybertronFromEncoderOptions struct {
	// PoolingStrategy specifies the pooling strategy for embedding calculation.
	PoolingStrategy int
}

// CybertronOptions represents options for the Cybertron embedder.
type CybertronOptions struct {
	CybertronFromEncoderOptions
	// ModelName is the name of the model (format: <org>/<model>).
	Model string
	// ModelsDir is the directory where the models are stored.
	ModelsDir string
	// HubAccessToken is the access token for the Hugging Face Hub.
	HubAccessToken string
}

// Cybertron represents an embedder powered by Cybertron.
type Cybertron struct {
	encoder textencoding.Interface
	opts    CybertronFromEncoderOptions
}

// NewCybertron creates a new instance of the Cybertron embedder.
func NewCybertron(optFns ...func(o *CybertronOptions)) (*Cybertron, error) {
	opts := CybertronOptions{
		Model:     "sentence-transformers/all-MiniLM-L6-v2",
		ModelsDir: "models",
		CybertronFromEncoderOptions: CybertronFromEncoderOptions{
			PoolingStrategy: int(bert.MeanPooling),
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	encoder, err := tasks.Load[textencoding.Interface](&tasks.Config{
		ModelsDir:      opts.ModelsDir,
		ModelName:      opts.Model,
		HubAccessToken: opts.HubAccessToken,
	})
	if err != nil {
		return nil, err
	}

	return NewCybertronFromEncoder(encoder, func(o *CybertronFromEncoderOptions) {
		o.PoolingStrategy = opts.PoolingStrategy
	})
}

// NewCybertronFromEncoder creates a new Cybertron embedder from an existing encoder.
func NewCybertronFromEncoder(encoder textencoding.Interface, optFns ...func(o *CybertronFromEncoderOptions)) (*Cybertron, error) {
	opts := CybertronFromEncoderOptions{
		PoolingStrategy: int(bert.MeanPooling),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Cybertron{
		encoder: encoder,
		opts:    opts,
	}, nil
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Cybertron) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		embedding, err := e.EmbedText(ctx, text)
		if err != nil {
			return nil, err
		}

		embeddings[i] = embedding
	}

	return embeddings, nil
}

// EmbedText embeds a single text and returns its embedding.
func (e *Cybertron) EmbedText(ctx context.Context, text string) ([]float32, error) {
	embedding, err := e.encoder.Encode(ctx, text, e.opts.PoolingStrategy)
	if err != nil {
		return nil, err
	}

	return embedding.Vector.Normalize2().Data().F32(), nil
}
