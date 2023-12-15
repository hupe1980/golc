package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestOpenAI(t *testing.T) {
	t.Run("EmbedDocuments", func(t *testing.T) {
		t.Run("Successful embedding", func(t *testing.T) {
			// Create a custom mock client.
			mockClient := &mockOpenAIClient{}

			// Create an instance of the OpenAI model with the custom mock client.
			openAIModel, err := NewOpenAIFromClient(mockClient)
			assert.NoError(t, err)

			// Define test inputs and expected outputs.
			texts := []string{"text1", "text2"}

			mockClient.response = openai.EmbeddingResponse{
				Data: []openai.Embedding{
					{
						Embedding: []float32{1.0, 2.0, 3.0},
					},
					{
						Embedding: []float32{4.0, 5.0, 6.0},
					},
				},
			}

			// Embed the documents.
			embeddings, err := openAIModel.BatchEmbedText(context.Background(), texts)

			// Assertions using testify/assert.
			assert.NoError(t, err, "Expected no error")
			assert.Len(t, embeddings, len(texts), "Expected the same number of embeddings as input texts")
		})

		t.Run("Test embedding error", func(t *testing.T) {
			// Create a custom mock client.
			mockClient := &mockOpenAIClient{}

			// Create an instance of the OpenAI model with the custom mock client.
			openAIModel, err := NewOpenAIFromClient(mockClient)
			assert.NoError(t, err)

			// Define test inputs.
			texts := []string{"text1", "text2"}

			// Configure the custom mock client to return an error.
			mockClient.err = errors.New("Test error")

			// Embed the documents.
			embeddings, err := openAIModel.BatchEmbedText(context.Background(), texts)

			// Assertions using testify/assert.
			assert.Error(t, err, "Expected an error")
			assert.Nil(t, embeddings, "Expected nil embeddings")
		})
	})

	t.Run("EmbedQuery", func(t *testing.T) {
		t.Run("Successful embedding", func(t *testing.T) {
			// Create a custom mock client.
			mockClient := &mockOpenAIClient{}

			// Create an instance of the OpenAI model with the custom mock client.
			openAIModel, err := NewOpenAIFromClient(mockClient)
			assert.NoError(t, err)

			mockClient.response = openai.EmbeddingResponse{
				Data: []openai.Embedding{
					{
						Embedding: []float32{1.0, 2.0, 3.0},
					},
				},
			}

			// Embed the documents.
			embeddings, err := openAIModel.EmbedText(context.Background(), "text1")

			// Assertions using testify/assert.
			assert.NoError(t, err, "Expected no error")
			assert.Len(t, embeddings, 3)
		})

		t.Run("Test embedding error", func(t *testing.T) {
			// Create a custom mock client.
			mockClient := &mockOpenAIClient{}

			// Create an instance of the OpenAI model with the custom mock client.
			openAIModel, err := NewOpenAIFromClient(mockClient)
			assert.NoError(t, err)

			// Configure the custom mock client to return an error.
			mockClient.err = errors.New("Test error")

			// Embed the documents.
			embeddings, err := openAIModel.EmbedText(context.Background(), "text1")

			// Assertions using testify/assert.
			assert.Error(t, err, "Expected an error")
			assert.Nil(t, embeddings, "Expected nil embeddings")
		})
	})
}

// mockOpenAIClient is a custom mock implementation of the OpenAI client interface.
type mockOpenAIClient struct {
	response openai.EmbeddingResponse
	err      error
}

// CreateEmbeddings mocks the CreateEmbeddings method of the OpenAI client.
func (m *mockOpenAIClient) CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error) {
	if m.err != nil {
		return openai.EmbeddingResponse{}, m.err
	}

	return m.response, nil
}
