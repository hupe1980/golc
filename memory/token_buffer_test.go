package memory

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/memory/chatmessagehistory"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/stretchr/testify/assert"
)

func TestConversationTokenBuffer(t *testing.T) {
	gpt2, err := tokenizer.NewGPT2()
	assert.NoError(t, err)

	cb := NewConversationTokenBuffer(gpt2)

	t.Run("MemoryKeys", func(t *testing.T) {
		expected := []string{"history"}
		variables := cb.MemoryKeys()
		assert.ElementsMatch(t, expected, variables)
	})

	t.Run("LoadMemoryVariables", func(t *testing.T) {
		inputs := map[string]any{}

		messages := schema.ChatMessages{
			schema.NewHumanChatMessage("Hello1"),
			schema.NewAIChatMessage("Hi there1"),
			schema.NewHumanChatMessage("Hello2"),
			schema.NewAIChatMessage("Hi there2"),
			schema.NewHumanChatMessage("Hello3"),
			schema.NewAIChatMessage("Hi there3"),
		}

		cb.opts.ChatMessageHistory = chatmessagehistory.NewInMemoryWithMessages(messages)

		t.Run("string - no shorten", func(t *testing.T) {
			cb.opts.ReturnMessages = false
			cb.opts.MaxTokenLimit = 2000

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, "Human: Hello1\nAI: Hi there1\nHuman: Hello2\nAI: Hi there2\nHuman: Hello3\nAI: Hi there3", vars["history"].(string))
		})

		t.Run("string - shorten", func(t *testing.T) {
			cb.opts.ReturnMessages = false
			cb.opts.MaxTokenLimit = 10

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, "Human: Hello3\nAI: Hi there3", vars["history"].(string))
		})

		t.Run("string - no history", func(t *testing.T) {
			cb.opts.ReturnMessages = false
			cb.opts.MaxTokenLimit = 0

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, "", vars["history"].(string))
		})

		t.Run("messages - no shorten", func(t *testing.T) {
			cb.opts.ReturnMessages = true
			cb.opts.MaxTokenLimit = 2000

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, 6, len(vars["history"].(schema.ChatMessages)))
		})

		t.Run("messages - shorten", func(t *testing.T) {
			cb.opts.ReturnMessages = true
			cb.opts.MaxTokenLimit = 10

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, 2, len(vars["history"].(schema.ChatMessages)))
		})

		t.Run("messages - no history", func(t *testing.T) {
			cb.opts.ReturnMessages = true
			cb.opts.MaxTokenLimit = 0

			vars, err := cb.LoadMemoryVariables(context.TODO(), inputs)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(vars["history"].(schema.ChatMessages)))
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
