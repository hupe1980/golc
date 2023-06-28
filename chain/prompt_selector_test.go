package chain

import (
	"testing"

	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestConditionalPromptSelector(t *testing.T) {
	t.Run("GetPrompt", func(t *testing.T) {
		defaultPrompt := prompt.NewTemplate("Default Prompt")
		llmPrompt := prompt.NewTemplate("LLM Prompt")
		chatModelPrompt := prompt.NewTemplate("ChatModel Prompt")

		conditional1 := Conditional{
			Condition: func(model schema.Model) bool {
				return IsLLM(model)
			},
			Prompt: llmPrompt,
		}

		conditional2 := Conditional{
			Condition: func(model schema.Model) bool {
				return IsChatModel(model)
			},
			Prompt: chatModelPrompt,
		}

		cps := ConditionalPromptSelector{
			DefaultPrompt: defaultPrompt,
			Conditionals:  []Conditional{conditional1, conditional2},
		}

		t.Run("LLM model should return LLM Prompt", func(t *testing.T) {
			llmModel := llm.NewFake("dummy")
			prompt := cps.GetPrompt(llmModel)

			text, _ := prompt.Format(nil)
			assert.Equal(t, "LLM Prompt", text)
		})

		t.Run("Chat model should return Chat Prompt", func(t *testing.T) {
			chatModel := chatmodel.NewFake("dummy")
			prompt := cps.GetPrompt(chatModel)

			text, _ := prompt.Format(nil)
			assert.Equal(t, "ChatModel Prompt", text)
		})
	})
}

func TestIsLLM(t *testing.T) {
	t.Run("LLM model should return true", func(t *testing.T) {
		llmModel := llm.NewFake("dummy")
		isLLM := IsLLM(llmModel)
		assert.True(t, isLLM)
	})

	t.Run("ChatModel should return false", func(t *testing.T) {
		chatModel := chatmodel.NewFake("dummy")
		isLLM := IsLLM(chatModel)
		assert.False(t, isLLM)
	})
}

func TestIsChatModel(t *testing.T) {
	t.Run("Chat model should return true", func(t *testing.T) {
		chatModel := chatmodel.NewFake("dummy")
		isChatModel := IsChatModel(chatModel)
		assert.True(t, isChatModel)
	})

	t.Run("LLM should return false", func(t *testing.T) {
		otherModel := llm.NewFake("dummy")
		isChatModel := IsChatModel(otherModel)
		assert.False(t, isChatModel)
	})
}
