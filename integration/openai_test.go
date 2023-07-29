package integration

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestToOpenAIChatCompletionMessages(t *testing.T) {
	messages := schema.ChatMessages{
		schema.NewAIChatMessage("Hello, how can I assist you?"),
		schema.NewHumanChatMessage("What is 1 times 1?"),
	}

	openAIMessages, err := ToOpenAIChatCompletionMessages(messages)
	require.NoError(t, err)
	require.Len(t, openAIMessages, 2)

	require.Equal(t, "assistant", openAIMessages[0].Role)
	require.Equal(t, "Hello, how can I assist you?", openAIMessages[0].Content)

	require.Equal(t, "user", openAIMessages[1].Role)
	require.Equal(t, "What is 1 times 1?", openAIMessages[1].Content)
}

// Test case for messageTypeToOpenAIRole function
func TestMessageTypeToOpenAIRole(t *testing.T) {
	assertRole, assertErr := messageTypeToOpenAIRole(schema.ChatMessageTypeAI)
	require.Equal(t, "assistant", assertRole)
	require.NoError(t, assertErr)

	unknownRole, unknownErr := messageTypeToOpenAIRole("unknown")
	require.Equal(t, "", unknownRole)
	require.EqualError(t, unknownErr, "unknown message type: unknown")
}
