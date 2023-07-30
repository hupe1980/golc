package agent

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestOpenAIFunctions(t *testing.T) {
	t.Run("TestPlan", func(t *testing.T) {
		agent, err := NewOpenAIFunctions(chatmodel.NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			var generation schema.Generation
			if len(messages) == 2 {
				generation = schema.Generation{
					Text: "text",
					Message: schema.NewAIChatMessage("text", func(o *schema.ChatMessageExtension) {
						o.FunctionCall = &schema.FunctionCall{
							Name:      "Mock",
							Arguments: `{"__arg1": "tool input"}`,
						}
					}),
				}
			} else {
				require.Len(t, messages, 4)
				require.Equal(t, "tool output", messages[3].Content())

				generation = schema.Generation{
					Text:    "finish text",
					Message: schema.NewAIChatMessage("finish text"),
				}
			}

			return &schema.ModelResult{
				Generations: []schema.Generation{generation},
				LLMOutput:   map[string]any{},
			}, nil
		}, func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{
				ToolRunFunc: func(ctx context.Context, input any) (string, error) {
					require.Equal(t, "tool input", input.(string))
					return "tool output", nil
				},
			},
		})
		require.NoError(t, err)

		// Create the inputs for the agent
		inputs := schema.ChainValues{
			"input": "User Input",
		}

		// Execute the agent's Plan method
		output, err := agent.Call(context.Background(), inputs)
		require.NoError(t, err)
		require.Equal(t, "finish text", output[agent.OutputKeys()[0]])
	})

	t.Run("TestPlanInvalidModel", func(t *testing.T) {
		_, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo"), []schema.Tool{
			&mockTool{},
		})
		require.Error(t, err)
		require.EqualError(t, err, "agent only supports OpenAI chatModels")
	})

	t.Run("TestPlanInvalidTool", func(t *testing.T) {
		_, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{ToolArgsType: struct {
				Channel chan int `json:"channel"` // chan cannot converted to json
			}{}},
		})
		require.Error(t, err)
		require.EqualError(t, err, "unsupported type chan from chan int")
	})
}
