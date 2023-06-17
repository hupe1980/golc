package memory

import (
	"testing"

	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestConversationBuffer(t *testing.T) {
	cb := NewConversationBuffer()

	t.Run("MemoryVariables", func(t *testing.T) {
		expected := []string{"history"}
		variables := cb.MemoryVariables()
		assert.ElementsMatch(t, expected, variables)
	})

	t.Run("LoadMemoryVariables", func(t *testing.T) {
		inputs := map[string]interface{}{}

		messages := []schema.ChatMessage{
			schema.NewHumanChatMessage("Hello"),
			schema.NewAIChatMessage("Hi there"),
		}

		cb.opts.ChatMessageHistory = chatmessagehistory.NewInMemoryWithMessages(messages)

		t.Run("ReturnMessages", func(t *testing.T) {
			// Test case with ReturnMessages set to true
			cb.opts.ReturnMessages = true
			expectedVariables := map[string]interface{}{
				"history": messages,
			}

			vars, err := cb.LoadMemoryVariables(inputs)
			assert.NoError(t, err)
			assert.Equal(t, len(expectedVariables), len(vars))
			for k, v := range vars {
				expected, ok := expectedVariables[k]
				assert.True(t, ok, "Unexpected memory variable '%s'", k)
				assert.Equal(t, expected, v, "Unexpected value for memory variable '%s'", k)
			}
		})

		t.Run("NoReturnMessages", func(t *testing.T) {
			// Test case with ReturnMessages set to false
			cb.opts.ReturnMessages = false
			expectedBuffer := "Human: Hello\nAI: Hi there"
			expectedVariables := map[string]interface{}{
				"history": expectedBuffer,
			}

			vars, err := cb.LoadMemoryVariables(inputs)
			assert.NoError(t, err)
			assert.Equal(t, len(expectedVariables), len(vars))
			for k, v := range vars {
				expected, ok := expectedVariables[k]
				assert.True(t, ok, "Unexpected memory variable '%s'", k)
				assert.Equal(t, expected, v, "Unexpected value for memory variable '%s'", k)
			}
		})
	})

	t.Run("SaveContext", func(t *testing.T) {
		inputs := map[string]interface{}{
			"input": "Hello",
		}
		outputs := map[string]interface{}{
			"output": "Hi there",
		}

		err := cb.SaveContext(inputs, outputs)
		assert.NoError(t, err)
	})
}
