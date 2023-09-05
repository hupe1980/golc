package chatmodel

import (
	"context"
	"testing"

	generativelanguagepb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
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
		mockClient.GenerateResponse = &generativelanguagepb.GenerateMessageResponse{
			Candidates: []*generativelanguagepb.Message{{
				Author:  "ai",
				Content: "World",
			}},
		}

		// Invoke the Generate method
		result, err := palm.Generate(context.Background(), schema.ChatMessages{
			schema.NewHumanChatMessage("Hello"),
		})

		// Assert the result and error
		assert.NoError(t, err)
		assert.Equal(t, "World", result.Generations[0].Message.Content())
	})

	t.Run("Type", func(t *testing.T) {
		// Create a Palm instance
		llm, err := NewPalm(&mockPalmClient{})
		assert.NoError(t, err)

		// Call the Type method
		typ := llm.Type()

		// Assert the result
		assert.Equal(t, "chatmodel.Palm", typ)
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
		})
		assert.NoError(t, err)

		// Call the InvocationParams method
		params := llm.InvocationParams()

		// Assert the result
		assert.Equal(t, float32(0.7), params["temperature"])
	})
}

// mockPalmClient is a mock implementation of the PalmClient interface for testing.
type mockPalmClient struct {
	GenerateResponse *generativelanguagepb.GenerateMessageResponse
	GenerateError    error
}

// GenerateMessage is a mock implementation of the GenerateMessage method.
func (m *mockPalmClient) GenerateMessage(ctx context.Context, req *generativelanguagepb.GenerateMessageRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateMessageResponse, error) {
	return m.GenerateResponse, m.GenerateError
}
