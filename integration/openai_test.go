package integration

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestToOpenAIChatCompletionMessages(t *testing.T) {
	messages := schema.ChatMessages{
		schema.NewAIChatMessage("Hello, how can I assist you?"),
		schema.NewHumanChatMessage("What is 1 times 1?"),
	}

	openAIMessages, err := ToOpenAIChatCompletionMessages(messages)
	assert.NoError(t, err)
	assert.Len(t, openAIMessages, 2)

	assert.Equal(t, "assistant", openAIMessages[0].Role)
	assert.Equal(t, "Hello, how can I assist you?", openAIMessages[0].Content)

	assert.Equal(t, "user", openAIMessages[1].Role)
	assert.Equal(t, "What is 1 times 1?", openAIMessages[1].Content)
}

// Test case for messageTypeToOpenAIRole function
func TestMessageTypeToOpenAIRole(t *testing.T) {
	assertRole, assertErr := messageTypeToOpenAIRole(schema.ChatMessageTypeAI)
	assert.Equal(t, "assistant", assertRole)
	assert.NoError(t, assertErr)

	unknownRole, unknownErr := messageTypeToOpenAIRole("unknown")
	assert.Equal(t, "", unknownRole)
	assert.EqualError(t, unknownErr, "unknown message type: unknown")
}
