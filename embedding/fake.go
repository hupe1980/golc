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

func (f *Fake) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))
	for i := range texts {
		embeddings[i] = f.getEmbedding()
	}

	return embeddings, nil
}

func (f *Fake) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	return f.getEmbedding(), nil
}

func (f *Fake) getEmbedding() []float64 {
	embedding := make([]float64, f.Size)
	for i := range embedding {
		embedding[i] = rand.NormFloat64()
	}

	return embedding
}
