package tokenizer

import (
	"context"
	"testing"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestGoogleGenAI(t *testing.T) {
	// Create a new instance of the GoogleGenAI model with the custom mock client.
	client := &mockGoogleGenAIClient{}

	// Create an instance of the GoogleGenAI tokenizer.
	GoogleGenAI := NewGoogleGenAI(client, "model")

	// Test GetNumTokens.
	t.Run("GetNumTokens", func(t *testing.T) {
		// Set the mock response for EmbedText.
		client.respCountTokens = &generativelanguagepb.CountTokensResponse{
			TotalTokens: 6,
		}
		client.errCount = nil

		// Test case with a sample input.
		text := "This is a sample text."
		numTokens, err := GoogleGenAI.GetNumTokens(context.TODO(), text)
		require.NoError(t, err)
		require.Equal(t, 6, int(numTokens))
	})

	// Test GetNumTokensFromMessage.
	t.Run("GetNumTokensFromMessage", func(t *testing.T) {
		// Set the mock response for EmbedText.
		client.respCountTokens = &generativelanguagepb.CountTokensResponse{
			TotalTokens: 27,
		}
		client.errCount = nil

		// Test case with sample chat messages.
		messages := schema.ChatMessages{
			schema.NewSystemChatMessage("Welcome to the chat!"),
			schema.NewHumanChatMessage("Hi, how are you?"),
			schema.NewSystemChatMessage("I'm doing well, thank you!"),
		}

		numTokens, err := GoogleGenAI.GetNumTokensFromMessage(context.TODO(), messages)
		require.NoError(t, err)
		require.Equal(t, 27, int(numTokens))
	})
}

// mockGoogleGenAIClient is a custom mock implementation of the GoogleGenAIClient interface.
type mockGoogleGenAIClient struct {
	respCountTokens *generativelanguagepb.CountTokensResponse
	errCount        error
}

// CountTokens mocks the CountTokens method of the GoogleGenAIClient interface.
func (m *mockGoogleGenAIClient) CountTokens(context.Context, *generativelanguagepb.CountTokensRequest, ...gax.CallOption) (*generativelanguagepb.CountTokensResponse, error) {
	if m.errCount != nil {
		return nil, m.errCount
	}

	return m.respCountTokens, nil
}
