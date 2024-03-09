package chatmodel

import (
	"context"
	"testing"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/core"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestCohereGenerate(t *testing.T) {
	// Create a new instance of the Cohere model with the mockCohereClient.
	mockClient := &mockCohereClient{}
	cohereModel, err := NewCohereFromClient(mockClient)
	assert.NoError(t, err)

	// Mock the Chat method to return a non-streamed response.
	mockClient.ChatFn = func(ctx context.Context, request *cohere.ChatRequest, opts ...core.RequestOption) (*cohere.NonStreamedChatResponse, error) {
		// Mock the response as needed for your test case.
		return &cohere.NonStreamedChatResponse{
			Text: "Mocked response",
		}, nil
	}

	t.Run("Generate", func(t *testing.T) {
		// Call the Generate method with your test case inputs.
		result, err := cohereModel.Generate(context.Background(), schema.ChatMessages{
			schema.NewHumanChatMessage("hello"),
		})
		assert.NoError(t, err)

		// Assert the expected result using testify assert.
		assert.NotNil(t, result)
		assert.Equal(t, "Mocked response", result.Generations[0].Text)
	})

	t.Run("no message", func(t *testing.T) {
		// Call the Generate method with your test case inputs.
		_, actualErr := cohereModel.Generate(context.Background(), schema.ChatMessages{})
		assert.ErrorContains(t, actualErr, "at least one message must be passed")
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "chatmodel.Cohere", cohereModel.Type())
	})

	// Test case for Callbacks method
	t.Run("Callbacks", func(t *testing.T) {
		assert.Equal(t, cohereModel.opts.CallbackOptions.Callbacks, cohereModel.Callbacks())
	})

	// Test case for Verbose method
	t.Run("Verbose", func(t *testing.T) {
		assert.Equal(t, cohereModel.opts.CallbackOptions.Verbose, cohereModel.Verbose())
	})

	// Test case for InvocationParams method
	t.Run("InvocationParams", func(t *testing.T) {
		// Call the InvocationParams method
		params := cohereModel.InvocationParams()

		// Assert the result
		assert.Equal(t, "command", params["model"])
		assert.Equal(t, 0.75, params["temperature"])
	})
}

// mockCohereClient is a mock implementation of the CohereClient interface.
type mockCohereClient struct {
	ChatFn       func(ctx context.Context, request *cohere.ChatRequest, opts ...core.RequestOption) (*cohere.NonStreamedChatResponse, error)
	ChatStreamFn func(ctx context.Context, request *cohere.ChatStreamRequest, opts ...core.RequestOption) (*core.Stream[cohere.StreamedChatResponse], error)
}

func (m *mockCohereClient) Chat(ctx context.Context, request *cohere.ChatRequest, opts ...core.RequestOption) (*cohere.NonStreamedChatResponse, error) {
	return m.ChatFn(ctx, request, opts...)
}

func (m *mockCohereClient) ChatStream(ctx context.Context, request *cohere.ChatStreamRequest, opts ...core.RequestOption) (*core.Stream[cohere.StreamedChatResponse], error) {
	return m.ChatStreamFn(ctx, request, opts...)
}
