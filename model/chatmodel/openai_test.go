package chatmodel

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestOpenAI_Generate(t *testing.T) {
	mockClient := &mockOpenAIClient{}
	openAI, err := NewOpenAIFromClient(mockClient)
	assert.NoError(t, err)

	// Test case for valid generation
	t.Run("ValidGeneration", func(t *testing.T) {
		ctx := context.Background()
		messages := schema.ChatMessages{
			schema.NewHumanChatMessage("Hello"),
			schema.NewAIChatMessage("Hi there"),
		}

		// Define the expected arguments and response for the mock client
		expectedRequest := openai.ChatCompletionRequest{
			Model:       openAI.opts.ModelName,
			Temperature: 1,
			TopP:        1,
			N:           0,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there"},
			},
			ToolChoice: "auto",
		}
		mockResponse := openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Role:    "assistant",
						Content: "Generated text",
					},
				},
			},
		}
		mockClient.createChatCompletionFn = func(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			assert.Equal(t, expectedRequest, request)
			return mockResponse, nil
		}

		result, err := openAI.Generate(ctx, messages)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Generations, 1)
		assert.Equal(t, "Generated text", result.Generations[0].Text)
	})

	// Test case for error during generation
	t.Run("GenerationError", func(t *testing.T) {
		ctx := context.Background()
		messages := schema.ChatMessages{
			schema.NewHumanChatMessage("Hello"),
			schema.NewAIChatMessage("Hi there"),
		}

		// Define the expected arguments and error for the mock client
		expectedRequest := openai.ChatCompletionRequest{
			Model:       openAI.opts.ModelName,
			Temperature: 1,
			TopP:        1,
			N:           0,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there"},
			},
			ToolChoice: "auto",
		}
		mockError := errors.New("generation error")
		mockClient.createChatCompletionFn = func(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
			assert.Equal(t, expectedRequest, request)
			return openai.ChatCompletionResponse{}, mockError
		}

		result, err := openAI.Generate(ctx, messages)
		assert.Error(t, err)
		assert.EqualError(t, errors.New("All attempts fail:\n#1: generation error"), err.Error())
		assert.Nil(t, result)
	})
	// Test case for Type method
	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "chatmodel.OpenAI", openAI.Type())
	})

	// Test case for Verbose method
	t.Run("Verbose", func(t *testing.T) {
		assert.Equal(t, openAI.opts.CallbackOptions.Verbose, openAI.Verbose())
	})

	// Test case for Callbacks method
	t.Run("Callbacks", func(t *testing.T) {
		assert.Equal(t, openAI.opts.CallbackOptions.Callbacks, openAI.Callbacks())
	})

	// Test case for InvocationParams method
	t.Run("InvocationParams", func(t *testing.T) {
		// Call the InvocationParams method
		params := openAI.InvocationParams()

		// Assert the result
		assert.Equal(t, "gpt-3.5-turbo", params["model_name"])
		assert.Equal(t, float32(1), params["temperature"])
	})
}

type mockOpenAIClient struct {
	createChatCompletionFn func(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (m *mockOpenAIClient) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if m.createChatCompletionFn != nil {
		return m.createChatCompletionFn(ctx, request)
	}

	return openai.ChatCompletionResponse{}, nil
}

func (m *mockOpenAIClient) CreateChatCompletionStream(ctx context.Context, request openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error) {
	return nil, nil
}

// Test case for openAIResponseToChatMessage function
func TestOpenAIResponseToChatMessage(t *testing.T) {
	aiMessage := openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: "AI response",
	}
	humanMessage := openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Human message",
	}

	// Assistant message
	assistantChatMessage := openAIResponseToChatMessage(aiMessage)
	assert.IsType(t, &schema.AIChatMessage{}, assistantChatMessage)
	assert.Equal(t, "AI response", assistantChatMessage.Content())

	// Human message
	humanChatMessage := openAIResponseToChatMessage(humanMessage)
	assert.IsType(t, &schema.HumanChatMessage{}, humanChatMessage)
	assert.Equal(t, "Human message", humanChatMessage.Content())

	// Unknown role message
	unknownChatMessage := openAIResponseToChatMessage(openai.ChatCompletionMessage{Content: "Unknown message", Role: "unknown"})
	assert.IsType(t, &schema.GenericChatMessage{}, unknownChatMessage)
	assert.Equal(t, "Unknown message", unknownChatMessage.Content())
}
