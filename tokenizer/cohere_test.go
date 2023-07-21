package tokenizer

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestCohere(t *testing.T) {
	// Create an instance of the Cohere tokenizer with a specific model.
	modelName := "model_name"
	cohere, err := NewCohere(modelName)
	require.NoError(t, err)

	// Test GetTokenIDs.
	t.Run("GetTokenIDs", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		ids, err := cohere.GetTokenIDs(text)
		require.NoError(t, err)
		require.ElementsMatch(t, []uint{1313, 329, 258, 7280, 2554, 47}, ids)
	})

	// Test GetNumTokens.
	t.Run("GetNumTokens", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		numTokens, err := cohere.GetNumTokens(text)
		require.NoError(t, err)
		require.Equal(t, 6, int(numTokens))
	})

	// Test GetNumTokensFromMessage.
	t.Run("GetNumTokensFromMessage", func(t *testing.T) {
		// Test case with sample chat messages.
		messages := schema.ChatMessages{
			schema.NewSystemChatMessage("Welcome to the chat!"),
			schema.NewHumanChatMessage("Hi, how are you?"),
			schema.NewSystemChatMessage("I'm doing well, thank you!"),
		}

		numTokens, err := cohere.GetNumTokensFromMessage(messages)
		require.NoError(t, err)
		require.Equal(t, 27, int(numTokens))
	})
}
