package llm

import (
	"context"
	"errors"
	"testing"

	"github.com/cohere-ai/cohere-go"
	"github.com/stretchr/testify/assert"
)

func TestCohere(t *testing.T) {
	ctx := context.Background()
	prompt := "Once upon a time"

	t.Run("Generate", func(t *testing.T) {
		// Create a mock client
		mockClient := &mockCohereClient{}

		// Create a Cohere instance with the mock client
		llm, err := NewCohereFromClient(mockClient)
		assert.NoError(t, err)

		t.Run("Successful generation", func(t *testing.T) {
			// Define the expected response from the mock client
			expectedResponse := &cohere.GenerateResponse{
				Generations: []cohere.Generation{{
					Text: "Once upon a time, there was a magical kingdom.",
				}},
			}

			// Mock the Generate method of the mock client
			mockClient.GenerateFunc = func(opts cohere.GenerateOptions) (*cohere.GenerateResponse, error) {
				return expectedResponse, nil
			}

			// Call the Generate method of the Cohere instance
			result, err := llm.Generate(ctx, prompt)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedResponse.Generations[0].Text, result.Generations[0].Text)
		})

		t.Run("Error in generation", func(t *testing.T) {
			// Define the error to be returned from the mock client
			returnedError := errors.New("generation failed")

			// Mock the Generate method of the mock client to return an error
			mockClient.GenerateFunc = func(opts cohere.GenerateOptions) (*cohere.GenerateResponse, error) {
				return nil, returnedError
			}

			// Call the Generate method of the Cohere instance
			result, err := llm.Generate(ctx, prompt)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Type", func(t *testing.T) {
		// Create a Cohere instance
		llm, err := NewCohereFromClient(&mockCohereClient{})
		assert.NoError(t, err)

		// Call the Type method
		typ := llm.Type()

		// Assert the result
		assert.Equal(t, "llm.Cohere", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		// Create a Cohere instance
		llm, err := NewCohereFromClient(&mockCohereClient{})
		assert.NoError(t, err)

		// Call the Verbose method
		verbose := llm.Verbose()

		// Assert the result
		assert.False(t, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		// Create a Cohere instance
		llm, err := NewCohereFromClient(&mockCohereClient{})
		assert.NoError(t, err)

		// Call the Callbacks method
		callbacks := llm.Callbacks()

		// Assert the result
		assert.Empty(t, callbacks)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Create a Cohere instance
		llm, err := NewCohereFromClient(&mockCohereClient{}, func(o *CohereOptions) {
			o.Model = "dummy"
			o.MaxTokens = 4711
		})
		assert.NoError(t, err)

		// Call the InvocationParams method
		params := llm.InvocationParams()

		// Assert the result
		assert.Equal(t, "dummy", params["model"])
		assert.Equal(t, uint(4711), params["max_tokens"])
	})
}

// mockCohereClient is a mock implementation of the CohereClient interface.
type mockCohereClient struct {
	GenerateFunc func(opts cohere.GenerateOptions) (*cohere.GenerateResponse, error)
}

// Generate is the mock implementation of the Generate method.
func (m *mockCohereClient) Generate(opts cohere.GenerateOptions) (*cohere.GenerateResponse, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(opts)
	}

	return nil, nil
}
