package chatmodel

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestConvertMessagesToAnthropicPrompt(t *testing.T) {
	t.Run("Empty input messages", func(t *testing.T) {
		emptyMessages := schema.ChatMessages{}
		emptyPrompt, emptyErr := convertMessagesToAnthropicPrompt(emptyMessages)
		assert.Equal(t, "", emptyPrompt)
		assert.Nil(t, emptyErr)
	})

	t.Run("Messages with a single system message", func(t *testing.T) {
		systemMessage := schema.NewSystemChatMessage("System message")
		messagesWithSystem := schema.ChatMessages{systemMessage}
		systemPrompt, systemErr := convertMessagesToAnthropicPrompt(messagesWithSystem)
		expectedSystemPrompt := "\n\nHuman: <admin>System message</admin>\n\nAssistant:"
		assert.Equal(t, expectedSystemPrompt, systemPrompt)
		assert.Nil(t, systemErr)
	})

	t.Run("Messages with a single AI message", func(t *testing.T) {
		aiMessage := schema.NewAIChatMessage("AI message")
		messagesWithAI := schema.ChatMessages{aiMessage}
		aiPrompt, aiErr := convertMessagesToAnthropicPrompt(messagesWithAI)
		expectedAIPrompt := "\n\nAssistant: AI message"
		assert.Equal(t, expectedAIPrompt, aiPrompt)
		assert.Nil(t, aiErr)
	})

	t.Run("Messages with a single human message", func(t *testing.T) {
		humanMessage := schema.NewHumanChatMessage("Human message")
		messagesWithHuman := schema.ChatMessages{humanMessage}
		humanPrompt, humanErr := convertMessagesToAnthropicPrompt(messagesWithHuman)
		expectedHumanPrompt := "\n\nHuman: Human message\n\nAssistant:"
		assert.Equal(t, expectedHumanPrompt, humanPrompt)
		assert.Nil(t, humanErr)
	})
}
