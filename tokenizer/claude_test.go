package tokenizer

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestClaude(t *testing.T) {
	// Create an instance of the Claude.
	claude, err := NewClaude()
	require.NoError(t, err)

	// Test GetTokenIDs.
	t.Run("GetTokenIDs", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		ids, err := claude.GetTokenIDs(context.TODO(), text)
		require.NoError(t, err)
		require.ElementsMatch(t, []uint{10545, 1800, 1320, 12110, 6840, 65}, ids)
	})

	// Test GetNumTokens.
	t.Run("GetNumTokens", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		numTokens, err := claude.GetNumTokens(context.TODO(), text)
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

		numTokens, err := claude.GetNumTokensFromMessage(context.TODO(), messages)
		require.NoError(t, err)
		require.Equal(t, 27, int(numTokens))
	})
}
