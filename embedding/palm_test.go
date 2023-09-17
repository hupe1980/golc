package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
)

func TestPalm(t *testing.T) {
	// Create a new instance of the Palm model with the custom mock client.
	client := &mockPalmClient{}
	palmModel := NewPalm(client)

	// Test cases
	t.Run("Test embedding of documents", func(t *testing.T) {
		// Define a list of texts to embed.
		texts := []string{"text1", "text2"}

		// Define a mock response for EmbedText.
		mockResponse := &generativelanguagepb.EmbedTextResponse{
			Embedding: &generativelanguagepb.Embedding{
				Value: []float32{1.0, 2.0, 3.0},
			},
		}

		// Set the mock response for EmbedText.
		client.respEmbedText = mockResponse
		client.errEmbedText = nil

		// Embed the documents.
		embeddings, err := palmModel.EmbedDocuments(context.Background(), texts)

		// Use assertions to check the results.
		assert.NoError(t, err)
		assert.NotNil(t, embeddings)
		assert.Len(t, embeddings, 2)
		assert.Len(t, embeddings[0], 3)
	})

	t.Run("Test embedding error", func(t *testing.T) {
		// Define a list of texts to embed.
		texts := []string{"text1"}

		// Set the mock error for EmbedText.
		client.respEmbedText = nil
		client.errEmbedText = errors.New("Test error")

		// Embed the documents.
		embeddings, err := palmModel.EmbedDocuments(context.Background(), texts)

		// Use assertions to check the error and embeddings.
		assert.Error(t, err)
		assert.Nil(t, embeddings)
	})

	t.Run("Test embedding of a single query", func(t *testing.T) {
		// Define a query text.
		query := "query text"

		// Define a mock response for EmbedText.
		mockResponse := &generativelanguagepb.EmbedTextResponse{
			Embedding: &generativelanguagepb.Embedding{
				Value: []float32{1.0, 2.0, 3.0},
			},
		}

		// Set the mock response for EmbedText.
		client.respEmbedText = mockResponse
		client.errEmbedText = nil

		// Embed the query.
		embedding, err := palmModel.EmbedQuery(context.Background(), query)

		// Use assertions to check the results.
		assert.NoError(t, err)
		assert.NotNil(t, embedding)
		assert.Len(t, embedding, 3)
	})

	t.Run("Test embedding error for query", func(t *testing.T) {
		// Define a query text.
		query := "query text"

		// Set the mock error for EmbedText.
		client.respEmbedText = nil
		client.errEmbedText = errors.New("Test error")

		// Embed the query.
		embedding, err := palmModel.EmbedQuery(context.Background(), query)

		// Use assertions to check the error and embedding.
		assert.Error(t, err)
		assert.Nil(t, embedding)
	})
}

// mockPalmClient is a custom mock implementation of the PalmClient interface.
type mockPalmClient struct {
	respEmbedText *generativelanguagepb.EmbedTextResponse
	errEmbedText  error
}

// EmbedText mocks the EmbedText method of the PalmClient interface.
func (m *mockPalmClient) EmbedText(ctx context.Context, req *generativelanguagepb.EmbedTextRequest, opts ...gax.CallOption) (*generativelanguagepb.EmbedTextResponse, error) {
	if m.errEmbedText != nil {
		return nil, m.errEmbedText
	}

	return m.respEmbedText, nil
}
