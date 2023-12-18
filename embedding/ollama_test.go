package embedding

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/integration/ollama"
	"github.com/stretchr/testify/assert"
)

func TestOllama(t *testing.T) {
	t.Run("EmbedText", func(t *testing.T) {
		client := &ollamaClientMock{}
		embedder := NewOllama(client)

		t.Run("Success", func(t *testing.T) {
			client.CreateEmbeddingFunc = func(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error) {
				return &ollama.EmbeddingResponse{Embedding: []float32{1.0, 2.0}}, nil
			}

			result, err := embedder.EmbedText(context.Background(), "text1")
			assert.NoError(t, err)
			assert.Equal(t, []float32{1.0, 2.0}, result)
		})

		t.Run("ErrorFromOllamaClient", func(t *testing.T) {
			expectedError := errors.New("error from OllamaClient")

			client.CreateEmbeddingFunc = func(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error) {
				return nil, expectedError
			}

			result, err := embedder.EmbedText(context.Background(), "text1")
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.EqualError(t, err, expectedError.Error())
		})
	})

	t.Run("BatchEmbedText", func(t *testing.T) {
		client := &ollamaClientMock{}
		embedder := NewOllama(client)

		t.Run("Success", func(t *testing.T) {
			texts := []string{"text1", "text2"}

			client.CreateEmbeddingFunc = func(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error) {
				return &ollama.EmbeddingResponse{Embedding: []float32{1.0, 2.0}}, nil
			}

			result, err := embedder.BatchEmbedText(context.Background(), texts)
			assert.NoError(t, err)
			assert.Equal(t, [][]float32{{1.0, 2.0}, {1.0, 2.0}}, result)
		})

		t.Run("ErrorFromOllamaClient", func(t *testing.T) {
			texts := []string{"text1", "text2"}
			expectedError := errors.New("error from OllamaClient")

			client.CreateEmbeddingFunc = func(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error) {
				return nil, expectedError
			}

			result, err := embedder.BatchEmbedText(context.Background(), texts)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.EqualError(t, err, expectedError.Error())
		})
	})
}

// ollamaClientMock is a custom mock implementation of the OllamaClient interface for testing purposes.
type ollamaClientMock struct {
	CreateEmbeddingFunc func(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error)
}

func (m *ollamaClientMock) CreateEmbedding(ctx context.Context, req *ollama.EmbeddingRequest) (*ollama.EmbeddingResponse, error) {
	if m.CreateEmbeddingFunc != nil {
		return m.CreateEmbeddingFunc(ctx, req)
	}

	panic("CreateEmbeddingFunc is not set in ollamaClientMock")
}
