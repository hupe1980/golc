package tokenizer

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestOpenAI(t *testing.T) {
	// Create an instance of the OpenAI tokenizer with a specific model name.
	modelName := "gpt-3.5-turbo" // Replace with your desired model name.
	openAI := NewOpenAI(modelName)

	// Test GetTokenIDs.
	t.Run("GetTokenIDs", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		ids, err := openAI.GetTokenIDs(text)
		require.NoError(t, err)
		require.ElementsMatch(t, []uint{2028, 374, 264, 6205, 1495, 13}, ids)
	})

	// Test GetNumTokens.
	t.Run("GetNumTokens", func(t *testing.T) {
		// Test case with a sample input.
		text := "This is a sample text."
		numTokens, err := openAI.GetNumTokens(text)
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

		numTokens, err := openAI.GetNumTokensFromMessage(messages)
		require.NoError(t, err)
		require.Equal(t, 25, int(numTokens))
	})
}
