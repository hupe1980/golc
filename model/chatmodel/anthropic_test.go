package chatmodel

import (
	"context"
	"fmt"
	"testing"

	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

// MockAnthropicClient is a mock implementation of the AnthropicClient interface for testing.
type MockAnthropicClient struct {
	createCompletionFn func(ctx context.Context, request *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error)
}

func (m *MockAnthropicClient) CreateCompletion(ctx context.Context, request *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
	return m.createCompletionFn(ctx, request)
}

func TestAnthropic(t *testing.T) {
	// Create a new instance of the Anthropic model with a mock client.
	client := &MockAnthropicClient{}

	// Initialize the Anthropic model with the mock client.
	anthropicModel, err := NewAnthropicFromClient(client)
	assert.NoError(t, err)

	t.Run("Generation", func(t *testing.T) {
		// Test case 1: Successful generation
		t.Run("Successful generation", func(t *testing.T) {
			// Mock the CreateCompletion method to return a valid response.
			client.createCompletionFn = func(ctx context.Context, request *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
				return &anthropic.CompletionResponse{
					Completion: "Hello, how can I help you?",
				}, nil
			}

			// Define chat messages
			chatMessages := []schema.ChatMessage{
				schema.NewAIChatMessage("Hi"),
				schema.NewHumanChatMessage("Can you help me?"),
			}

			// Generate text
			result, err := anthropicModel.Generate(context.Background(), chatMessages)
			assert.NoError(t, err, "Expected no error")
			assert.NotNil(t, result, "Expected non-nil result")
			assert.Len(t, result.Generations, 1, "Expected 1 generation")
			assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
		})

		// Test case 2: Anthropic API error
		t.Run("Anthropic API error", func(t *testing.T) {
			// Mock the CreateCompletion method to return an error response.
			client.createCompletionFn = func(ctx context.Context, request *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error) {
				return nil, fmt.Errorf("Anthropic API error")
			}

			// Define chat messages
			chatMessages := []schema.ChatMessage{
				schema.NewAIChatMessage("Hi"),
				schema.NewHumanChatMessage("Can you help me?"),
			}

			// Generate text
			result, err := anthropicModel.Generate(context.Background(), chatMessages)
			assert.Error(t, err, "Expected an error")
			assert.Nil(t, result, "Expected nil result")
		})
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "chatmodel.Anthropic", anthropicModel.Type())
	})

	t.Run("Callbacks", func(t *testing.T) {
		assert.Equal(t, anthropicModel.opts.CallbackOptions.Callbacks, anthropicModel.Callbacks())
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Call the InvocationParams method
		params := anthropicModel.InvocationParams()

		// Assert the result
		assert.Equal(t, "claude-v1", params["model_name"])
		assert.Equal(t, float32(0.5), params["temperature"])
	})
}

func TestConvertMessagesToAnthropicPrompt(t *testing.T) {
	t.Run("Empty input messages", func(t *testing.T) {
		emptyMessages := schema.ChatMessages{}
		emptyPrompt, emptyErr := convertMessagesToAnthropicPrompt(emptyMessages)
		assert.Equal(t, "", emptyPrompt)
		assert.Nil(t, emptyErr)
	})

	t.Run("Messages with a single system message", func(t *testing.T) {
		systemMessage := schema.NewSystemChatMessage("System message")
		messagesWithSystem := schema.ChatMessages{systemMessage}
		systemPrompt, systemErr := convertMessagesToAnthropicPrompt(messagesWithSystem)
		expectedSystemPrompt := "\n\nHuman: <admin>System message</admin>\n\nAssistant:"
		assert.Equal(t, expectedSystemPrompt, systemPrompt)
		assert.Nil(t, systemErr)
	})

	t.Run("Messages with a single AI message", func(t *testing.T) {
		aiMessage := schema.NewAIChatMessage("AI message")
		messagesWithAI := schema.ChatMessages{aiMessage}
		aiPrompt, aiErr := convertMessagesToAnthropicPrompt(messagesWithAI)
		expectedAIPrompt := "\n\nAssistant: AI message"
		assert.Equal(t, expectedAIPrompt, aiPrompt)
		assert.Nil(t, aiErr)
	})

	t.Run("Messages with a single human message", func(t *testing.T) {
		humanMessage := schema.NewHumanChatMessage("Human message")
		messagesWithHuman := schema.ChatMessages{humanMessage}
		humanPrompt, humanErr := convertMessagesToAnthropicPrompt(messagesWithHuman)
		expectedHumanPrompt := "\n\nHuman: Human message\n\nAssistant:"
		assert.Equal(t, expectedHumanPrompt, humanPrompt)
		assert.Nil(t, humanErr)
	})
}
