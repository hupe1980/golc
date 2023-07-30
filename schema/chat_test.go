package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatMessageToMap(t *testing.T) {
	humanMsg := NewHumanChatMessage("Hello, I am a human.")
	aiMsg := NewAIChatMessage("Hello, I am an AI.")
	funcMsg := NewFunctionChatMessage("foo", "bar")

	humanMap := ChatMessageToMap(humanMsg)
	aiMap := ChatMessageToMap(aiMsg)
	funcMap := ChatMessageToMap(funcMsg)

	require.Equal(t, "human", humanMap["type"])
	require.Equal(t, "Hello, I am a human.", humanMap["content"])

	require.Equal(t, "ai", aiMap["type"])
	require.Equal(t, "Hello, I am an AI.", aiMap["content"])

	require.Equal(t, "function", funcMap["type"])
	require.Equal(t, "foo", funcMap["name"])
	require.Equal(t, "bar", funcMap["content"])
}

func TestMapToChatMessage(t *testing.T) {
	humanMap := map[string]string{
		"type":    "human",
		"content": "Hello, I am a human.",
	}
	aiMap := map[string]string{
		"type":    "ai",
		"content": "Hello, I am an AI.",
	}

	humanMsg, err := MapToChatMessage(humanMap)
	require.NoError(t, err)
	require.IsType(t, &HumanChatMessage{}, humanMsg)
	require.Equal(t, "Hello, I am a human.", humanMsg.Content())

	aiMsg, err := MapToChatMessage(aiMap)
	require.NoError(t, err)
	require.IsType(t, &AIChatMessage{}, aiMsg)
	require.Equal(t, "Hello, I am an AI.", aiMsg.Content())
}

func TestStringifyChatMessages(t *testing.T) {
	chatMessages := ChatMessages{
		NewHumanChatMessage("Hello, I am a human."),
		NewAIChatMessage("Hello, I am an AI."),
		NewSystemChatMessage("System message."),
		NewGenericChatMessage("Generic message.", "role"),
		NewFunctionChatMessage("function", "Function call message."),
	}

	formatted, err := chatMessages.Format()
	require.NoError(t, err)
	require.Contains(t, formatted, "Human: Hello, I am a human.")
	require.Contains(t, formatted, "AI: Hello, I am an AI.")
	require.Contains(t, formatted, "System: System message.")
	require.Contains(t, formatted, "role: Generic message.")
	require.Contains(t, formatted, "Function: Function call message.")
}
