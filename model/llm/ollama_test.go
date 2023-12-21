package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hupe1980/golc/integration/ollama"
)

func TestOllama(t *testing.T) {
	t.Parallel()

	t.Run("Generate", func(t *testing.T) {
		t.Parallel()

		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			mockClient := &mockOllamaClient{
				GenerateFunc: func(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationResponse, error) {
					assert.Equal(t, "llama2", req.Model)
					assert.Equal(t, "Hello", req.Prompt)

					return &ollama.GenerationResponse{
						Response: "I can help you with that.",
					}, nil
				},
			}

			ollamaModel, err := NewOllama(mockClient)
			assert.NoError(t, err)

			// Run the model
			result, err := ollamaModel.Generate(context.Background(), "Hello")
			assert.NoError(t, err)

			// Check the result
			assert.Len(t, result.Generations, 1)
			assert.Equal(t, "I can help you with that.", result.Generations[0].Text)
		})

		t.Run("Error", func(t *testing.T) {
			t.Parallel()

			mockClient := &mockOllamaClient{
				GenerateFunc: func(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationResponse, error) {
					return nil, errors.New("error generating chat")
				},
			}

			ollamaModel, err := NewOllama(mockClient)
			assert.NoError(t, err)

			result, err := ollamaModel.Generate(context.Background(), "Hello")
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Type", func(t *testing.T) {
		t.Parallel()

		mockClient := &mockOllamaClient{}
		ollamaModel, err := NewOllama(mockClient)
		assert.NoError(t, err)

		assert.Equal(t, "llm.Ollama", ollamaModel.Type())
	})

	t.Run("Callbacks", func(t *testing.T) {
		t.Parallel()

		mockClient := &mockOllamaClient{}
		ollamaModel, err := NewOllama(mockClient)
		assert.NoError(t, err)

		assert.Equal(t, ollamaModel.opts.CallbackOptions.Callbacks, ollamaModel.Callbacks())
	})

	t.Run("InvocationParams", func(t *testing.T) {
		t.Parallel()

		mockClient := &mockOllamaClient{}
		ollamaModel, err := NewOllama(mockClient)
		assert.NoError(t, err)

		params := ollamaModel.InvocationParams()
		assert.NotNil(t, params)
		assert.Equal(t, float32(0.7), params["temperature"])
		assert.Equal(t, 256, params["max_tokens"])
	})
}

// mockOllamaClient is a mock implementation of the llm.OllamaClient interface.
type mockOllamaClient struct {
	GenerateFunc func(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationResponse, error)
}

// CreateGeneration is the mock implementation of the CreateGeneration method for mockOllamaClient.
func (m *mockOllamaClient) CreateGeneration(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationResponse, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, req)
	}

	return nil, errors.New("GenerateChatFunc not implemented")
}

// CreateGenerationStream is the mock implementation of the CreateGenerationStream method for mockOllamaClient.
func (m *mockOllamaClient) CreateGenerationStream(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationStream, error) {
	return nil, nil
}
