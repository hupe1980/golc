package prompt

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestStringPromptValue(t *testing.T) {
	t.Run("String method", func(t *testing.T) {
		// Test cases for the String method.
		stringValueTests := []struct {
			name     string
			input    StringPromptValue
			expected string
		}{
			{
				name:     "StringPromptValue with text",
				input:    StringPromptValue("Hello, World!"),
				expected: "Hello, World!",
			},
			{
				name:     "Empty StringPromptValue",
				input:    StringPromptValue(""),
				expected: "",
			},
		}

		for _, test := range stringValueTests {
			t.Run(test.name, func(t *testing.T) {
				result := test.input.String()
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("Messages method", func(t *testing.T) {
		// Test cases for the Messages method.
		messagesTests := []struct {
			name     string
			input    StringPromptValue
			expected schema.ChatMessages
		}{
			{
				name:  "StringPromptValue with text",
				input: StringPromptValue("Hello, World!"),
				expected: schema.ChatMessages{
					schema.NewHumanChatMessage("Hello, World!"),
				},
			},
			{
				name:  "Empty StringPromptValue",
				input: StringPromptValue(""),
				expected: schema.ChatMessages{
					schema.NewHumanChatMessage(""),
				},
			},
		}

		for _, test := range messagesTests {
			t.Run(test.name, func(t *testing.T) {
				result := test.input.Messages()
				assert.Equal(t, test.expected, result)
			})
		}
	})
}
