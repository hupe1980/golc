package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/integration/ai21"
	"github.com/stretchr/testify/assert"
)

// MockAI21Client is a custom mock implementation of the AI21Client interface.
type MockAI21Client struct {
	CreateCompletionFunc func(ctx context.Context, model string, req *ai21.CompleteRequest) (*ai21.CompleteResponse, error)
}

// CreateCompletion mocks the CreateCompletion method of AI21Client.
func (m *MockAI21Client) CreateCompletion(ctx context.Context, model string, req *ai21.CompleteRequest) (*ai21.CompleteResponse, error) {
	if m.CreateCompletionFunc != nil {
		return m.CreateCompletionFunc(ctx, model, req)
	}

	return nil, errors.New("CreateCompletionFunc not implemented")
}

func TestAI21(t *testing.T) {
	// Initialize the AI21 client with a mock client
	mockClient := &MockAI21Client{}
	llm, err := NewAI21FromClient(mockClient, func(o *AI21Options) {
		o.Model = "j2-mid"
		o.Temperature = 0.7
		o.MaxTokens = 256
	})
	assert.NoError(t, err)

	t.Run("Generate_Success", func(t *testing.T) {
		// Set up expected values
		expectedPrompt := "Test prompt"
		expectedResponse := &ai21.CompleteResponse{
			Completions: []ai21.Completion{
				{
					Data: ai21.Data{
						Text: "Generated text",
					},
				},
			},
		}

		// Implement the CreateCompletionFunc for the mock client
		mockClient.CreateCompletionFunc = func(ctx context.Context, model string, req *ai21.CompleteRequest) (*ai21.CompleteResponse, error) {
			assert.Equal(t, expectedPrompt, req.Prompt) // Assert that the prompt is as expected
			return expectedResponse, nil
		}

		// Call the Generate method
		result, err := llm.Generate(context.Background(), expectedPrompt)

		// Assert the result and error
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResponse.Completions[0].Data.Text, result.Generations[0].Text)
	})

	t.Run("Generate_Error", func(t *testing.T) {
		// Set up expected values
		expectedPrompt := "Test prompt"
		expectedError := errors.New("AI21 error")

		// Implement the CreateCompletionFunc for the mock client to return an error
		mockClient.CreateCompletionFunc = func(ctx context.Context, model string, req *ai21.CompleteRequest) (*ai21.CompleteResponse, error) {
			assert.Equal(t, expectedPrompt, req.Prompt) // Assert that the prompt is as expected
			return nil, expectedError
		}

		// Call the Generate method
		result, err := llm.Generate(context.Background(), expectedPrompt)

		// Assert the error and result
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}
