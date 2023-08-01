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

func TestStructuredOutput(t *testing.T) {
	t.Run("TestStructuredOutput", func(t *testing.T) {
		// Create a dummy chat model for testing
		chatModel := chatmodel.NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			return &schema.ModelResult{
				Generations: []schema.Generation{{
					Text: "",
					Message: schema.NewAIChatMessage("", func(o *schema.ChatMessageExtension) {
						o.FunctionCall = &schema.FunctionCall{
							Name:      "Person",
							Arguments: `{"name": "Max", "age": 21}`,
						}
					}),
				}},
				LLMOutput: map[string]any{},
			}, nil
		})

		// Create a dummy prompt template for testing
		promptTemplate := prompt.NewChatTemplate([]prompt.MessageTemplate{
			prompt.NewHumanMessageTemplate("{{.input}}"),
		})

		// Create a dummy output candidate
		type person struct {
			Name    string `json:"name" description:"The person's name"`
			Age     int    `json:"age" description:"The person's age"`
			FavFood string `json:"fav_food,omitempty" description:"The person's favorite food"`
		}

		// Create a new StructuredOutput chain
		structuredOutputChain, err := NewStructuredOutput(chatModel, promptTemplate, []OutputCandidate{
			{
				Name:        "Person",
				Description: "Identifying information about a person",
				Data:        &person{},
			},
		})
		require.NoError(t, err)

		// Prepare the input values for the chain
		inputs := schema.ChainValues{
			"input": "Max is 21",
		}

		// Call the ChatModel chain with the inputs
		outputs, err := golc.Call(context.Background(), structuredOutputChain, inputs)
		require.NoError(t, err)

		// Check the output key in the result
		require.Contains(t, outputs, structuredOutputChain.OutputKeys()[0])

		// Check the output value type
		p, ok := outputs[structuredOutputChain.OutputKeys()[0]].(*person)
		require.True(t, ok)
		require.Equal(t, "Max", p.Name)
		require.Equal(t, 21, p.Age)
		require.Equal(t, "", p.FavFood)
	})
}
