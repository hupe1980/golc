package embedding

import (
	"context"
	"math/rand"
)

type Fake struct {
	Size int
}

func NewFake(size int) *Fake {
	return &Fake{Size: size}
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *Fake) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))
	for i := range texts {
		embeddings[i] = e.getEmbedding()
	}

	return embeddings, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *Fake) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	return e.getEmbedding(), nil
}

func (e *Fake) getEmbedding() []float64 {
	embedding := make([]float64, e.Size)
	for i := range embedding {
		embedding[i] = rand.NormFloat64()
	}

	return embedding
}
