package embedding

import (
	"context"
	"errors"
	"testing"

	cohere "github.com/cohere-ai/cohere-go/v2"
	core "github.com/cohere-ai/cohere-go/v2/core"
	"github.com/stretchr/testify/assert"
)

func TestCohere(t *testing.T) {
	t.Run("EmbedDocuments", func(t *testing.T) {
		t.Run("Successful embedding of documents", func(t *testing.T) {
			// Create a new instance of the Cohere model with a mock client.
			client := &mockCohereClient{
				response: &cohere.EmbedResponse{
					EmbeddingsByType: &cohere.EmbedByTypeResponse{
						Embeddings: &cohere.EmbedByTypeResponseEmbeddings{
							Float: [][]float64{
								{1.0, 2.0, 3.0},
								{4.0, 5.0, 6.0},
							},
						},
					},
				},
			}

			// Initialize the Cohere model with the mock client.
			cohereModel := NewCohereFromClient(client)

			// Define a list of texts to embed.
			texts := []string{"text1", "text2"}

			// Embed the documents.
			embeddings, err := cohereModel.BatchEmbedText(context.Background(), texts)
			assert.NoError(t, err, "Expected no error")
			assert.NotNil(t, embeddings, "Expected non-nil embeddings")
			assert.Len(t, embeddings, 2, "Expected 2 embeddings")
			assert.Len(t, embeddings[0], 3, "Expected 3 values in the embedding")
		})
	})

	t.Run("EmbedQuery", func(t *testing.T) {
		t.Run("Successful embedding of a single query", func(t *testing.T) {
			// Create a new instance of the Cohere model with a mock client.
			client := &mockCohereClient{
				response: &cohere.EmbedResponse{
					EmbeddingsByType: &cohere.EmbedByTypeResponse{
						Embeddings: &cohere.EmbedByTypeResponseEmbeddings{
							Float: [][]float64{
								{1.0, 2.0, 3.0},
							},
						},
					},
				},
			}

			// Initialize the Cohere model with the mock client.
			cohereModel := NewCohereFromClient(client)

			// Define a query text.
			query := "query text"

			// Embed the query.
			embedding, err := cohereModel.EmbedText(context.Background(), query)
			assert.NoError(t, err, "Expected no error")
			assert.NotNil(t, embedding, "Expected non-nil embedding")
			assert.Len(t, embedding, 3, "Expected 3 values in the embedding")
		})

		// Test case: Embedding error
		t.Run("Embedding error", func(t *testing.T) {
			// Create a new instance of the Cohere model with a mock client.
			client := &mockCohereClient{
				err: errors.New("Embedding error"),
			}

			// Initialize the Cohere model with the mock client.
			cohereModel := NewCohereFromClient(client)
			// // Mock the Embed method to return an error.
			// client.embedFn = func(opts cohere.EmbedOptions) (*cohere.EmbedResponse, error) {
			// 	return nil, errors.New("Embedding error")
			// }

			// Define a query text.
			query := "query text"

			// Embed the query.
			embedding, err := cohereModel.EmbedText(context.Background(), query)
			assert.Error(t, err, "Expected an error")
			assert.Nil(t, embedding, "Expected nil embedding")
		})
	})
}

// mockCohereClient is a mock implementation of the CohereClient interface for testing.
type mockCohereClient struct {
	response *cohere.EmbedResponse
	err      error
}

func (m *mockCohereClient) Embed(ctx context.Context, request *cohere.EmbedRequest, opts ...core.RequestOption) (*cohere.EmbedResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.response, nil
}
