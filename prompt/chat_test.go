package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hupe1980/golc/schema"
)

func TestChatTemplateWrapper(t *testing.T) {
	// Create some sample chat templates
	chatTemplate1 := NewSystemMessageTemplate("Welcome, {{.Name}}!")
	chatTemplate2 := NewAIMessageTemplate("Hello, I'm an AI.")
	chatTemplate3 := NewHumanMessageTemplate("How can I help you, {{.Name}}?")

	chatTemplates1 := NewChatTemplate([]MessageTemplate{chatTemplate1})
	chatTemplates2 := NewChatTemplate([]MessageTemplate{chatTemplate2, chatTemplate3})

	// Create the chat template wrapper
	chatTemplateWrapper := NewChatTemplateWrapper(chatTemplates1, chatTemplates2)

	// Define the input values
	values := map[string]interface{}{
		"Name": "John",
	}

	// Run the test cases
	t.Run("FormatPrompt", func(t *testing.T) {
		expectedMessages := schema.ChatMessages{
			schema.NewSystemChatMessage("Welcome, John!"),
			schema.NewAIChatMessage("Hello, I'm an AI."),
			schema.NewHumanChatMessage("How can I help you, John?"),
		}

		// Call the FormatPrompt method
		promptValue, err := chatTemplateWrapper.FormatPrompt(values)

		// Check the result
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, promptValue.Messages())
	})

	t.Run("Format", func(t *testing.T) {
		expectedMessages := schema.ChatMessages{
			schema.NewSystemChatMessage("Welcome, John!"),
			schema.NewAIChatMessage("Hello, I'm an AI."),
			schema.NewHumanChatMessage("How can I help you, John?"),
		}

		// Call the Format method
		messages, err := chatTemplateWrapper.Format(values)

		// Check the result
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)
	})
}

func TestMessagesPlaceholder(t *testing.T) {
	// Create a chat template placeholder
	placeholder := NewMessagesPlaceholder("Messages")

	// Define the input values
	values := map[string]interface{}{
		"Messages": schema.ChatMessages{
			schema.NewSystemChatMessage("Hello"),
			schema.NewHumanChatMessage("How are you?"),
			schema.NewAIChatMessage("I'm fine, thank you!"),
		},
	}

	// Run the test case
	t.Run("Format", func(t *testing.T) {
		expectedMessages, _ := values["Messages"].(schema.ChatMessages)

		// Call the Format method
		messages, err := placeholder.Format(values)

		// Check the result
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)
	})
}
