package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestChatModel(t *testing.T) {
	t.Run("TestChatModelCall", func(t *testing.T) {
		// Create a dummy chat model for testing
		chatModel := chatmodel.NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			text := messages[0].Content()

			return &schema.ModelResult{
				Generations: []schema.Generation{{
					Text:    text,
					Message: schema.NewAIChatMessage(text),
				}},
				LLMOutput: map[string]any{},
			}, nil
		})

		// Create a dummy prompt template for testing
		promptTemplate := prompt.NewChatTemplate([]prompt.MessageTemplate{
			prompt.NewHumanMessageTemplate("Hello {{.input}}"),
		})

		// Create a new ChatModel chain
		chatModelChain, err := NewChatModel(chatModel, promptTemplate)
		require.NoError(t, err)

		// Prepare the input values for the chain
		inputs := schema.ChainValues{
			"input": "World",
		}

		// Call the ChatModel chain with the inputs
		outputs, err := golc.Call(context.Background(), chatModelChain, inputs)
		require.NoError(t, err)

		// Check the output key in the result
		require.Contains(t, outputs, chatModelChain.OutputKeys()[0])

		// Check the output value type
		msg, ok := outputs[chatModelChain.OutputKeys()[0]].(*schema.AIChatMessage)
		require.True(t, ok)
		require.Equal(t, "Hello World", msg.Content())
	})
}
