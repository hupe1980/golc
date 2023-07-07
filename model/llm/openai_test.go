package llm

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestOpenAI_Generate(t *testing.T) {
	// Create a mock OpenAIClient
	mockClient := &mockOpenAIClient{}

	// Create an OpenAI instance with the mock client
	openAI, err := NewOpenAIFromClient(mockClient)
	assert.NoError(t, err)

	t.Run("SuccessfulCompletion", func(t *testing.T) {
		// Mock the completion stream and error
		mockClient.CompletionStream = nil
		mockClient.CompletionStreamErr = nil

		// Mock the completion response and error
		mockClient.CompletionResponse = openai.CompletionResponse{
			Choices: []openai.CompletionChoice{{
				Text:         "World",
				FinishReason: "stop",
			}},
			Usage: openai.Usage{
				PromptTokens:     10,
				CompletionTokens: 10,
				TotalTokens:      20,
			},
		}
		mockClient.CompletionResponseErr = nil

		// Expected result
		expectedResult := &schema.ModelResult{
			Generations: []schema.Generation{{
				Text: "World",
				Info: map[string]any{
					"FinishReason": "stop",
					"LogProbs":     openai.LogprobResult{},
				},
			}},
			LLMOutput: map[string]any{
				"ModelName": "text-davinci-002",
				"TokenUsage": map[string]int{
					"CompletionTokens": 10,
					"PromptTokens":     10,
					"TotalTokens":      20,
				},
			},
		}

		// Invoke the Generate method
		result, err := openAI.Generate(context.Background(), "Hello")

		// Assert the result and error
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})
}

// mockOpenAIClient is a mock implementation of the OpenAIClient interface for testing.
type mockOpenAIClient struct {
	CompletionStream      *openai.CompletionStream
	CompletionStreamErr   error
	CompletionResponse    openai.CompletionResponse
	CompletionResponseErr error
}

// CreateCompletionStream is a mock implementation of the CreateCompletionStream method.
func (m *mockOpenAIClient) CreateCompletionStream(ctx context.Context, request openai.CompletionRequest) (stream *openai.CompletionStream, err error) {
	return m.CompletionStream, m.CompletionStreamErr
}

// CreateCompletion is a mock implementation of the CreateCompletion method.
func (m *mockOpenAIClient) CreateCompletion(ctx context.Context, request openai.CompletionRequest) (response openai.CompletionResponse, err error) {
	return m.CompletionResponse, m.CompletionResponseErr
}
