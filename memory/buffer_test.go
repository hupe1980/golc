package memory

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestConversationBuffer(t *testing.T) {
	cb := NewConversationBuffer()

	t.Run("MemoryKeys", func(t *testing.T) {
		expected := []string{"history"}
		variables := cb.MemoryKeys()
		assert.ElementsMatch(t, expected, variables)
	})

	t.Run("LoadMemoryVariables", func(t *testing.T) {
		inputs := map[string]interface{}{}

		messages := schema.ChatMessages{
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

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
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

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, len(expectedVariables), len(vars))
			for k, v := range vars {
				expected, ok := expectedVariables[k]
				assert.True(t, ok, "Unexpected memory variable '%s'", k)
				assert.Equal(t, expected, v, "Unexpected value for memory variable '%s'", k)
			}
		})

		t.Run("Sliding Window", func(t *testing.T) {
			sw := NewConversationBuffer()

			messages := schema.ChatMessages{
				schema.NewHumanChatMessage("Hello1"),
				schema.NewAIChatMessage("Hi there1"),
				schema.NewHumanChatMessage("Hello2"),
				schema.NewAIChatMessage("Hi there2"),
				schema.NewHumanChatMessage("Hello3"),
				schema.NewAIChatMessage("Hi there3"),
			}

			sw.opts.ChatMessageHistory = chatmessagehistory.NewInMemoryWithMessages(messages)
			sw.opts.ReturnMessages = true

			t.Run("K=0", func(t *testing.T) {
				sw.opts.K = 0

				vars, err := sw.LoadMemoryVariables(context.TODO(), inputs)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(vars["history"].(schema.ChatMessages)))
			})

			t.Run("K=1 (-> 2 Messages)", func(t *testing.T) {
				sw.opts.K = 1

				vars, err := sw.LoadMemoryVariables(context.TODO(), inputs)
				assert.NoError(t, err)

				messages, _ := vars["history"].(schema.ChatMessages)
				assert.Equal(t, 2, len(messages))

				assert.Equal(t, "Hello3", messages[0].Text())
				assert.Equal(t, "Hi there3", messages[1].Text())
			})
		})
	})

	t.Run("SaveContext", func(t *testing.T) {
		inputs := map[string]interface{}{
			"input": "Hello",
		}
		outputs := map[string]interface{}{
			"output": "Hi there",
		}

		err := cb.SaveContext(context.TODO(), inputs, outputs)
		assert.NoError(t, err)
	})
}
