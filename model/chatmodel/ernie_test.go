package chatmodel

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/integration/ernie"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestErnie(t *testing.T) {
	// Create a new instance of the Ernie model with a mock client.
	client := &mockErnieClient{}

	// Initialize the Ernie model with the mock client.
	ernieModel, err := NewErnieFromClient(client)
	assert.NoError(t, err)

	t.Run("Generation", func(t *testing.T) {
		// Test case 1: Successful generation
		t.Run("Successful generation", func(t *testing.T) {
			// Mock the CreateChatCompletion method to return a valid response.
			client.createChatCompletionFn = func(ctx context.Context, model string, request *ernie.ChatCompletionRequest) (*ernie.ChatCompletionResponse, error) {
				return &ernie.ChatCompletionResponse{
					Result:    "Hello, how can I help you?",
					ErrorCode: 0,
				}, nil
			}

			// Define chat messages
			chatMessages := []schema.ChatMessage{
				schema.NewAIChatMessage("Hi"),
				schema.NewHumanChatMessage("Can you help me?"),
			}

			// Generate text
			result, err := ernieModel.Generate(context.Background(), chatMessages)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result.Generations, 1)
			assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text)
			assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Message.Content())
		})

		// Test case 2: Ernie API error
		t.Run("Ernie API error", func(t *testing.T) {
			// Mock the CreateChatCompletion method to return an error response.
			client.createChatCompletionFn = func(ctx context.Context, model string, request *ernie.ChatCompletionRequest) (*ernie.ChatCompletionResponse, error) {
				return &ernie.ChatCompletionResponse{
					ErrorCode: 123,
				}, nil
			}

			// Define chat messages
			chatMessages := []schema.ChatMessage{
				schema.NewAIChatMessage("Hi"),
				schema.NewHumanChatMessage("Can you help me?"),
			}

			// Generate text
			result, err := ernieModel.Generate(context.Background(), chatMessages)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})
}

// mockErnieClient is a mock implementation of the ErnieClient interface for testing.
type mockErnieClient struct {
	createChatCompletionFn func(ctx context.Context, model string, request *ernie.ChatCompletionRequest) (*ernie.ChatCompletionResponse, error)
}

func (m *mockErnieClient) CreateChatCompletion(ctx context.Context, model string, request *ernie.ChatCompletionRequest) (*ernie.ChatCompletionResponse, error) {
	return m.createChatCompletionFn(ctx, model, request)
}
