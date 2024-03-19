package vectorstore

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	t.Run("SaveAndLoad", func(t *testing.T) {
		originalData := []InMemoryItem{
			{Content: "item1", Vector: []float32{1.0, 2.0, 3.0}, Metadata: map[string]any{"key1": "value1"}},
			{Content: "item2", Vector: []float32{4.0, 5.0, 6.0}, Metadata: map[string]any{"key2": "value2"}},
		}

		// Create an InMemory instance with the original data
		vsOriginal := &InMemory{data: originalData}

		// Serialize the original data
		var buf bytes.Buffer
		err := vsOriginal.Save(&buf)
		require.NoError(t, err, "Failed to save data")

		// Create a new InMemory instance
		vsLoaded := &InMemory{}

		// Load the serialized data
		err = vsLoaded.Load(&buf)
		require.NoError(t, err, "Failed to load data")

		// Check if the loaded data matches the original data
		assert.Equal(t, originalData, vsLoaded.data, "Loaded data does not match original data")
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
