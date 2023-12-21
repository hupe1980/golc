package chatmodel

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestGoogleGenAI(t *testing.T) {
	mockClient := &mockGoogleGenAIClient{}
	model, err := NewGoogleGenAI(mockClient)
	assert.NoError(t, err)

	t.Run("Generate_Success", func(t *testing.T) {
		mockClient.GenerateContentFn = func(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateContentResponse, error) {
			// Implement your custom behavior here, e.g., return a predefined response
			return &generativelanguagepb.GenerateContentResponse{
				Candidates: []*generativelanguagepb.Candidate{{
					Content: &generativelanguagepb.Content{
						Parts: []*generativelanguagepb.Part{{Data: &generativelanguagepb.Part_Text{
							Text: "Generated text",
						}}},
					},
				}},
			}, nil
		}

		// Define chat messages
		chatMessages := []schema.ChatMessage{
			schema.NewHumanChatMessage("Can you help me?"),
		}

		result, err := model.Generate(context.Background(), chatMessages)
		assert.NoError(t, err)
		assert.Equal(t, "Generated text", result.Generations[0].Text)
		assert.Equal(t, "Generated text", result.Generations[0].Message.Content())
	})

	t.Run("Generate_Error", func(t *testing.T) {
		mockClient.GenerateContentFn = func(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateContentResponse, error) {
			// Implement your custom behavior here, e.g., return a predefined response
			return nil, fmt.Errorf("google genai error")
		}

		// Define chat messages
		chatMessages := []schema.ChatMessage{
			schema.NewHumanChatMessage("Can you help me?"),
		}

		_, err := model.Generate(context.Background(), chatMessages)
		assert.ErrorContains(t, err, "google genai error")
	})

	// Test the Type method
	t.Run("Type", func(t *testing.T) {
		expectedType := "chatmodel.GoogleGenAI"
		assert.Equal(t, expectedType, model.Type())
	})

	// Test the Verbose method
	t.Run("Verbose", func(t *testing.T) {
		assert.False(t, model.Verbose())
	})

	// Test the Callbacks method
	t.Run("Callbacks", func(t *testing.T) {
		callbacks := model.Callbacks()
		assert.Empty(t, callbacks)
	})

	// Test the InvocationParams method
	t.Run("InvocationParams", func(t *testing.T) {
		invocationParams := model.InvocationParams()
		assert.Equal(t, "models/gemini-pro", invocationParams["model_name"])
		assert.Equal(t, int32(1), invocationParams["candidate_count"])
	})
}

// mockGoogleGenAIClient is a custom mock implementation of the GoogleGenAIClient interface.
type mockGoogleGenAIClient struct {
	GenerateContentFn       func(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateContentResponse, error)
	StreamGenerateContentFn func(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (generativelanguagepb.GenerativeService_StreamGenerateContentClient, error)
	CountTokensFn           func(ctx context.Context, req *generativelanguagepb.CountTokensRequest, opts ...gax.CallOption) (*generativelanguagepb.CountTokensResponse, error)
}

// GenerateContent is a mocked method for the GenerateContent function.
func (m *mockGoogleGenAIClient) GenerateContent(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateContentResponse, error) {
	if m.GenerateContentFn != nil {
		return m.GenerateContentFn(ctx, req, opts...)
	}

	return nil, errors.New("GenerateContent not implemented in the mock")
}

// StreamGenerateContent is a mocked method for the StreamGenerateContent function.
func (m *mockGoogleGenAIClient) StreamGenerateContent(ctx context.Context, req *generativelanguagepb.GenerateContentRequest, opts ...gax.CallOption) (generativelanguagepb.GenerativeService_StreamGenerateContentClient, error) {
	if m.StreamGenerateContentFn != nil {
		return m.StreamGenerateContentFn(ctx, req, opts...)
	}

	return nil, errors.New("StreamGenerateContent not implemented in the mock")
}

// CountTokens is a mocked method for the CountTokens function.
func (m *mockGoogleGenAIClient) CountTokens(ctx context.Context, req *generativelanguagepb.CountTokensRequest, opts ...gax.CallOption) (*generativelanguagepb.CountTokensResponse, error) {
	if m.CountTokensFn != nil {
		return m.CountTokensFn(ctx, req, opts...)
	}

	return nil, errors.New("CountTokens not implemented in the mock")
}
