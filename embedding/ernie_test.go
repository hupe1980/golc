package embedding

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/integration/ernie"
	"github.com/stretchr/testify/assert"
)

func TestEmbedding(t *testing.T) {
	t.Run("EmbedDocuments", func(t *testing.T) {
		// Create an instance of Ernie with the mock client.
		ernieClient := &mockErnieClient{
			Response: &ernie.EmbeddingResponse{
				ID:     "fakeID",
				Object: "fakeObject",
				Data: []struct {
					Object    string    `json:"object"`
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Object:    "text1",
						Embedding: []float32{1.0, 2.0, 3.0},
						Index:     0,
					},
					{
						Object:    "text2",
						Embedding: []float32{4.0, 5.0, 6.0},
						Index:     1,
					},
					{
						Object:    "text3",
						Embedding: []float32{7.0, 8.0, 9.0},
						Index:     2,
					},
				},
			},
		}
		ernieEmbed := NewErnieFromClient(ernieClient)

		ctx := context.Background()
		texts := []string{"text1", "text2", "text3"}

		// Test embedding of documents.
		embeddings, err := ernieEmbed.BatchEmbedText(ctx, texts)

		assert.NoError(t, err, "Error embedding documents")
		assert.Len(t, embeddings, len(texts), "Unexpected number of embeddings")

		assert.ElementsMatch(t, []float32{1.0, 2.0, 3.0}, embeddings[0])
		assert.ElementsMatch(t, []float32{4.0, 5.0, 6.0}, embeddings[1])
		assert.ElementsMatch(t, []float32{7.0, 8.0, 9.0}, embeddings[2])
	})

	t.Run("EmbedQuery", func(t *testing.T) {
		// Create an instance of Ernie with the mock client.
		ernieClient := &mockErnieClient{
			Response: &ernie.EmbeddingResponse{
				ID:     "fakeID",
				Object: "fakeObject",
				Data: []struct {
					Object    string    `json:"object"`
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Object:    "fakeEmbedding",
						Embedding: []float32{1.0, 2.0, 3.0},
						Index:     0,
					},
				},
			},
		}
		ernieEmbed := NewErnieFromClient(ernieClient)

		ctx := context.Background()
		query := "queryText"

		// Test embedding of a query.
		embedding, err := ernieEmbed.EmbedText(ctx, query)

		assert.NoError(t, err, "Error embedding query")
		assert.Len(t, embedding, 3, "Unexpected embedding dimensions")
		expected := []float32{1.0, 2.0, 3.0} // Mocked embedding values
		assert.ElementsMatch(t, expected, embedding, "Embedding values do not match for the query")
	})
}

// mockErnieClient is a mock implementation of ErnieClient for testing purposes.
type mockErnieClient struct {
	Response *ernie.EmbeddingResponse
}

func (m *mockErnieClient) CreateEmbedding(ctx context.Context, model string, request ernie.EmbeddingRequest) (*ernie.EmbeddingResponse, error) {
	// Mock implementation, return fake embeddings.
	return m.Response, nil
}
