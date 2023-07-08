package llm

import (
	"context"
	"testing"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestVertexAI_Generate(t *testing.T) {
	// Create a mock VertexAIClient
	mockClient := &mockVertexAIClient{}

	// Create a VertexAI instance with the mock client
	vertexAI, err := NewVertexAI(mockClient, "dummy")
	assert.NoError(t, err)

	// Run the test case
	t.Run("SuccessfulGeneration", func(t *testing.T) {
		prediction, err := structpb.NewValue(map[string]any{
			"content": "World",
		})
		assert.NoError(t, err)

		mockClient.PredictResponse = &aiplatformpb.PredictResponse{
			Predictions: []*structpb.Value{prediction},
		}

		// Invoke the Generate method
		result, err := vertexAI.Generate(context.Background(), "Hello")

		// Assert the result and error
		assert.NoError(t, err)
		assert.Equal(t, "World", result.Generations[0].Text)
	})

	t.Run("Type", func(t *testing.T) {
		// Create a VertexAI instance
		llm, err := NewVertexAI(&mockVertexAIClient{}, "dummy")
		assert.NoError(t, err)

		// Call the Type method
		typ := llm.Type()

		// Assert the result
		assert.Equal(t, "llm.VertexAI", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		// Create a VertexAI instance
		llm, err := NewVertexAI(&mockVertexAIClient{}, "dummy")
		assert.NoError(t, err)

		// Call the Verbose method
		verbose := llm.Verbose()

		// Assert the result
		assert.False(t, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		// Create a VertexAI instance
		llm, err := NewVertexAI(&mockVertexAIClient{}, "dummy")
		assert.NoError(t, err)

		// Call the Callbacks method
		callbacks := llm.Callbacks()

		// Assert the result
		assert.Empty(t, callbacks)
	})

	t.Run("InvocationParams", func(t *testing.T) {
		// Create a VertexAI instance
		llm, err := NewVertexAI(&mockVertexAIClient{}, "dummy", func(o *VertexAIOptions) {
			o.Temperatur = 0.7
			o.MaxOutputTokens = 4711
		})
		assert.NoError(t, err)

		// Call the InvocationParams method
		params := llm.InvocationParams()

		// Assert the result
		assert.Equal(t, float32(0.7), params["temperatur"])
		assert.Equal(t, 4711, params["max_output_tokens"])
	})
}

// mockVertexAIClient is a mock implementation of the VertexAIClient interface for testing.
type mockVertexAIClient struct {
	PredictResponse *aiplatformpb.PredictResponse
	PredictError    error
}

// Predict is a mock implementation of the Predict method.
func (m *mockVertexAIClient) Predict(ctx context.Context, req *aiplatformpb.PredictRequest, opts ...gax.CallOption) (*aiplatformpb.PredictResponse, error) {
	return m.PredictResponse, m.PredictError
}
