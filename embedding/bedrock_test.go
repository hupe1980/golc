package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/assert"
)

func TestBedrock(t *testing.T) {
	t.Run("TestEmbedDocuments", func(t *testing.T) {
		t.Run("Successful embedding of documents", func(t *testing.T) {
			// Create an instance of the Bedrock struct with a mock client.
			client := &mockBedrockRuntimeClient{
				response: &bedrockruntime.InvokeModelOutput{
					Body: []byte(`{"embedding": [1.0, 2.0, 3.0]}`),
				},
			}
			embedder := NewBedrock(client)

			// Define a list of texts to embed.
			texts := []string{"text1", "text2"}

			// Embed the documents.
			embeddings, err := embedder.BatchEmbedText(context.Background(), texts)

			// Add your assertions using testify
			assert.NoError(t, err, "Expected no error")
			assert.NotNil(t, embeddings, "Expected non-nil embeddings")
			assert.Len(t, embeddings, 2, "Expected 2 embeddings")
			assert.Len(t, embeddings[0], 3, "Expected 3 values in the embedding")
		})
	})

	t.Run("TestEmbedQuery", func(t *testing.T) {
		t.Run("Successful embedding of a single query", func(t *testing.T) {
			// Create an instance of the Bedrock struct with a mock client.
			client := &mockBedrockRuntimeClient{
				response: &bedrockruntime.InvokeModelOutput{
					Body: []byte(`{"embedding": [1.0, 2.0, 3.0]}`),
				},
			}
			embedder := NewBedrock(client)

			// Define a text.
			text := "text"

			// Embed the text.
			embedding, err := embedder.EmbedText(context.Background(), text)

			// Add your assertions using testify
			assert.NoError(t, err, "Expected no error")
			assert.NotNil(t, embedding, "Expected non-nil embedding")
			assert.Len(t, embedding, 3, "Expected 3 values in the embedding")
		})

		t.Run("Embedding error", func(t *testing.T) {
			// Create an instance of the Bedrock struct with a mock client.
			client := &mockBedrockRuntimeClient{
				err: errors.New("Embedding error"),
			}
			embedder := NewBedrock(client)

			// Define a text.
			text := "text"

			// Embed the text.
			embedding, err := embedder.EmbedText(context.Background(), text)

			// Add your assertions using testify
			assert.Error(t, err, "Expected an error")
			assert.Nil(t, embedding, "Expected nil embedding")
		})
	})
}

// mockBedrockRuntimeClient is a mock implementation of BedrockRuntimeClient for testing.
type mockBedrockRuntimeClient struct {
	response *bedrockruntime.InvokeModelOutput
	err      error
}

func (m *mockBedrockRuntimeClient) InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.response, nil
}
