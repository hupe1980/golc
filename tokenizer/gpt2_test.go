package tokenizer

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestGPT2(t *testing.T) {
	// Create an instance of the GPT2 tokenizer.
	gpt2, err := NewGPT2()
	require.NoError(t, err)

	// Test GetTokenIDs.
	t.Run("GetTokenIDs", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		ids, err := gpt2.GetTokenIDs(text)
		require.NoError(t, err)
		require.ElementsMatch(t, []uint{1212, 318, 257, 6291, 2420, 13}, ids)
	})

	// Test GetNumTokens.
	t.Run("GetNumTokens", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		numTokens, err := gpt2.GetNumTokens(text)
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

		numTokens, err := gpt2.GetNumTokensFromMessage(messages)
		require.NoError(t, err)
		require.Equal(t, 27, int(numTokens))
	})
}
