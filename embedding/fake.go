package embedding

import (
	"context"
	"math/rand"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Fake satisfies the Embedder interface.
var _ schema.Embedder = (*Fake)(nil)

type Fake struct {
	Size int
}

func NewFake(size int) *Fake {
	return &Fake{Size: size}
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Fake) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = e.getEmbedding()
	}

	return embeddings, nil
}

// EmbedText embeds a single text and returns its embedding.
func (e *Fake) EmbedText(ctx context.Context, text string) ([]float32, error) {
	return e.getEmbedding(), nil
}

func (e *Fake) getEmbedding() []float32 {
	embedding := make([]float32, e.Size)
	for i := range embedding {
		embedding[i] = float32(rand.NormFloat64())
	}

	return embedding
}
