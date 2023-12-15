package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
)

func TestGoogleGenAI(t *testing.T) {
	// Create a new instance of the GoogleGenAI model with the custom mock client.
	client := &mockGoogleGenAIClient{}

	googleGenAIModel := NewGoogleGenAI(client)

	// Test cases
	t.Run("Test embedding of documents", func(t *testing.T) {
		// Define a list of texts to embed.
		texts := []string{"text1", "text2"}

		// Define a mock response for EmbedText.
		mockResponse := &generativelanguagepb.BatchEmbedContentsResponse{
			Embeddings: []*generativelanguagepb.ContentEmbedding{
				{
					Values: []float32{1.0, 2.0, 3.0},
				}, {
					Values: []float32{4.0, 5.0, 6.0},
				}},
		}

		// Set the mock response for EmbedText.
		client.respEmbedBatchContents = mockResponse
		client.errEmbed = nil

		// Embed the documents.
		embeddings, err := googleGenAIModel.BatchEmbedText(context.Background(), texts)

		// Use assertions to check the results.
		assert.NoError(t, err)
		assert.NotNil(t, embeddings)
		assert.Len(t, embeddings, 2)
		assert.Len(t, embeddings[0], 3)
		assert.Equal(t, float32(1.0), embeddings[0][0])
		assert.Len(t, embeddings[1], 3)
		assert.Equal(t, float32(4.0), embeddings[1][0])
	})

	t.Run("Test embedding error", func(t *testing.T) {
		// Define a list of texts to embed.
		texts := []string{"text1"}

		// Set the mock error for EmbedText.
		client.respEmbedBatchContents = nil
		client.errEmbed = errors.New("Test error")

		// Embed the documents.
		embeddings, err := googleGenAIModel.BatchEmbedText(context.Background(), texts)

		// Use assertions to check the error and embeddings.
		assert.Error(t, err)
		assert.Nil(t, embeddings)
	})

	t.Run("Test embedding of a single query", func(t *testing.T) {
		// Define a query text.
		query := "query text"

		// Define a mock response for EmbedText.
		mockResponse := &generativelanguagepb.EmbedContentResponse{
			Embedding: &generativelanguagepb.ContentEmbedding{
				Values: []float32{1.0, 2.0, 3.0},
			},
		}

		// Set the mock response for EmbedText.
		client.respEmbedContent = mockResponse
		client.errEmbed = nil

		// Embed the query.
		embedding, err := googleGenAIModel.EmbedText(context.Background(), query)

		// Use assertions to check the results.
		assert.NoError(t, err)
		assert.NotNil(t, embedding)
		assert.Len(t, embedding, 3)
	})

	t.Run("Test embedding error for query", func(t *testing.T) {
		// Define a query text.
		query := "query text"

		// Set the mock error for EmbedText.
		client.respEmbedContent = nil
		client.errEmbed = errors.New("Test error")

		// Embed the query.
		embedding, err := googleGenAIModel.EmbedText(context.Background(), query)

		// Use assertions to check the error and embedding.
		assert.Error(t, err)
		assert.Nil(t, embedding)
	})
}

// mockGoogleGenAIClient is a custom mock implementation of the GoogleGenAIClient interface.
type mockGoogleGenAIClient struct {
	respEmbedContent       *generativelanguagepb.EmbedContentResponse
	respEmbedBatchContents *generativelanguagepb.BatchEmbedContentsResponse
	errEmbed               error
}

// EmbedContent mocks the EmbedContent method of the GoogleGenAIClient interface.
func (m *mockGoogleGenAIClient) EmbedContent(context.Context, *generativelanguagepb.EmbedContentRequest, ...gax.CallOption) (*generativelanguagepb.EmbedContentResponse, error) {
	if m.errEmbed != nil {
		return nil, m.errEmbed
	}

	return m.respEmbedContent, nil
}

// BatchEmbedContents mocks the BatchEmbedContents method of the GoogleGenAIClient interface.
func (m *mockGoogleGenAIClient) BatchEmbedContents(context.Context, *generativelanguagepb.BatchEmbedContentsRequest, ...gax.CallOption) (*generativelanguagepb.BatchEmbedContentsResponse, error) {
	if m.errEmbed != nil {
		return nil, m.errEmbed
	}

	return m.respEmbedBatchContents, nil
}
