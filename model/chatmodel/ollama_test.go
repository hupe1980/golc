package chatmodel

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/schema"
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
				GenerateChatFunc: func(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatResponse, error) {
					assert.Equal(t, "llama2", req.Model)
					assert.Len(t, req.Messages, 2)
					assert.Equal(t, "user", req.Messages[0].Role)
					assert.Equal(t, "Hello", req.Messages[0].Content)
					assert.Equal(t, "assistant", req.Messages[1].Role)
					assert.Equal(t, "How can I help you?", req.Messages[1].Content)

					return &ollama.ChatResponse{
						Message: &ollama.Message{
							Role:    "assistant",
							Content: "I can help you with that.",
						},
					}, nil
				},
			}

			ollamaModel, err := NewOllama(mockClient)
			assert.NoError(t, err)

			// Simulate chat messages
			messages := []schema.ChatMessage{
				schema.NewHumanChatMessage("Hello"),
				schema.NewAIChatMessage("How can I help you?"),
			}

			// Run the model
			result, err := ollamaModel.Generate(context.Background(), messages)
			assert.NoError(t, err)

			// Check the result
			assert.Len(t, result.Generations, 1)
			assert.Equal(t, "I can help you with that.", result.Generations[0].Text)
		})

		t.Run("Error", func(t *testing.T) {
			t.Parallel()

			mockClient := &mockOllamaClient{
				GenerateChatFunc: func(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatResponse, error) {
					return nil, errors.New("error generating chat")
				},
			}

			ollamaModel, err := NewOllama(mockClient)
			assert.NoError(t, err)

			messages := []schema.ChatMessage{
				schema.NewHumanChatMessage("Hello"),
				schema.NewAIChatMessage("How can I help you?"),
			}

			result, err := ollamaModel.Generate(context.Background(), messages)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Type", func(t *testing.T) {
		t.Parallel()

		mockClient := &mockOllamaClient{}
		ollamaModel, err := NewOllama(mockClient)
		assert.NoError(t, err)

		assert.Equal(t, "chatmodel.Ollama", ollamaModel.Type())
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

// mockOllamaClient is a mock implementation of the chatmodel.OllamaClient interface.
type mockOllamaClient struct {
	GenerateChatFunc func(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatResponse, error)
}

// CreateChat is the mock implementation of the CreateChat method for mockOllamaClient.
func (m *mockOllamaClient) CreateChat(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatResponse, error) {
	if m.GenerateChatFunc != nil {
		return m.GenerateChatFunc(ctx, req)
	}

	return nil, errors.New("GenerateChatFunc not implemented")
}

// CreateChatStream is the mock implementation of the CreateChatStream method for mockOllamaClient.
func (m *mockOllamaClient) CreateChatStream(ctx context.Context, req *ollama.ChatRequest) (*ollama.ChatStream, error) {
	return nil, nil
}
