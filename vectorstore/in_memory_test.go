package vectorstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hupe1980/golc/schema"
)

func TestInMemory(t *testing.T) {
	// Setup
	embedder := &mockEmbedder{}
	vs := NewInMemory(embedder)

	// Test AddDocuments method
	t.Run("AddDocuments", func(t *testing.T) {
		// Given
		documents := []schema.Document{
			{PageContent: "document1"},
			{PageContent: "document2"},
			{PageContent: "document3"},
		}

		// When
		err := vs.AddDocuments(context.Background(), documents)

		// Then
		assert.NoError(t, err)
		assert.Len(t, vs.Data(), 3)
	})

	// Test SimilaritySearch method
	t.Run("SimilaritySearch", func(t *testing.T) {
		// Given
		query := "query"
		expectedDocuments := []schema.Document{
			{PageContent: "document1"},
			{PageContent: "document2"},
			{PageContent: "document3"},
		}

		// When
		documents, err := vs.SimilaritySearch(context.Background(), query)

		// Then
		assert.NoError(t, err)
		assert.Len(t, documents, 3)

		for i, doc := range documents {
			assert.Equal(t, expectedDocuments[i].PageContent, doc.PageContent)
		}
	})
}

// mockEmbedder implements the schema.Embedder interface for testing purposes.
type mockEmbedder struct{}

func (m *mockEmbedder) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	// Mock implementation for batch embedding text
	return [][]float32{
		{1.0, 2.0, 3.0},
		{2.0, 3.0, 4.0},
		{3.0, 4.0, 5.0},
	}, nil
}

func (m *mockEmbedder) EmbedText(ctx context.Context, text string) ([]float32, error) {
	// Mock implementation for embedding text
	return []float32{1.0, 2.0, 3.0}, nil
}
