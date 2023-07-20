package llm

import (
	"context"
	"testing"

	generativelanguagepb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestPalm(t *testing.T) {
	// Create a mock PalmClient
	mockClient := &mockPalmClient{}

	// Create a Palm instance with the mock client
	palm, err := NewPalm(mockClient)
	assert.NoError(t, err)

	// Run the test case
	t.Run("SuccessfulGeneration", func(t *testing.T) {
		mockClient.GenerateResponse = &generativelanguagepb.GenerateTextResponse{
			Candidates: []*generativelanguagepb.TextCompletion{{
				Output: "World",
			}},
		}

		// Invoke the Generate method
		result, err := palm.Generate(context.Background(), "Hello")

		// Assert the result and error
		assert.NoError(t, err)
		assert.Equal(t, "World", result.Generations[0].Text)
	})

	t.Run("Type", func(t *testing.T) {
		// Create a Palm instance
		llm, err := NewPalm(&mockPalmClient{})
		assert.NoError(t, err)

		// Call the Type method
		typ := llm.Type()

		// Assert the result
		assert.Equal(t, "llm.Palm", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		// Create a Palm instance
		llm, err := NewPalm(&mockPalmClient{})
		assert.NoError(t, err)

		// Call the Verbose method
		verbose := llm.Verbose()

		// Assert the result
		assert.False(t, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		// Create a Palm instance
		llm, err := NewPalm(&mockPalmClient{})
		assert.NoError(t, err)

		// Call the Callbacks method
		callbacks := llm.Callbacks()

		// Assert the result
		assert.Empty(t, callbacks)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Create a Palm instance
		llm, err := NewPalm(&mockPalmClient{}, func(o *PalmOptions) {
			o.Temperature = 0.7
			o.MaxOutputTokens = 4711
		})
		assert.NoError(t, err)

		// Call the InvocationParams method
		params := llm.InvocationParams()

		// Assert the result
		assert.Equal(t, float32(0.7), params["temperature"])
		assert.Equal(t, int32(4711), params["max_output_tokens"])
	})
}

// mockPalmClient is a mock implementation of the PalmClient interface for testing.
type mockPalmClient struct {
	GenerateResponse *generativelanguagepb.GenerateTextResponse
	GenerateError    error
}

// GenerateText is a mock implementation of the GenerateText method.
func (m *mockPalmClient) GenerateText(ctx context.Context, req *generativelanguagepb.GenerateTextRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateTextResponse, error) {
	return m.GenerateResponse, m.GenerateError
}
