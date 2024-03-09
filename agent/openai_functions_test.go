package agent

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestOpenAIFunctions(t *testing.T) {
	t.Parallel()

	t.Run("TestPlan", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
			var generation schema.Generation

			if len(messages) == 2 {
				assert.Equal(t, "user Input", messages[1].Content())

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
				assert.Len(t, messages, 4)
				assert.Equal(t, "tool output", messages[3].Content())

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
					assert.Equal(t, "tool input", input.(string))
					return "tool output", nil
				},
			},
		})
		assert.NoError(t, err)

		// Create the inputs for the agent
		inputs := schema.ChainValues{
			"input": "user Input",
		}

		// Execute the agent's Plan method
		output, err := agent.Call(context.Background(), inputs)
		assert.NoError(t, err)
		assert.Equal(t, "finish text", output[agent.OutputKeys()[0]])
	})

	t.Run("TestPlanInvalidModel", func(t *testing.T) {
		t.Parallel()

		_, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo"), []schema.Tool{
			&mockTool{},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "agent only supports OpenAI chatModels")
	})

	t.Run("TestPlanInvalidTool", func(t *testing.T) {
		t.Parallel()

		_, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{ToolArgsType: struct {
				Channel chan int `json:"channel"` // chan cannot converted to json
			}{}},
		})
		assert.Error(t, err)
		assert.EqualError(t, err, "unsupported type chan from chan int")
	})

	t.Run("InputKeys", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{},
		})
		assert.NoError(t, err)

		keys := agent.InputKeys()
		assert.ElementsMatch(t, keys, []string{"input"})
	})

	t.Run("OutputKeys", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{},
		})
		assert.NoError(t, err)

		keys := agent.OutputKeys()
		assert.ElementsMatch(t, keys, []string{"output"})
	})

	t.Run("Type", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{},
		})
		assert.NoError(t, err)

		typ := agent.Type()
		assert.Equal(t, "OpenAIFunctions", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{},
		})
		assert.NoError(t, err)

		verbose := agent.Verbose()

		assert.Equal(t, agent.opts.CallbackOptions.Verbose, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		t.Parallel()

		agent, err := NewOpenAIFunctions(chatmodel.NewSimpleFake("foo", func(o *chatmodel.FakeOptions) {
			o.ChatModelType = "chatmodel.OpenAI"
		}), []schema.Tool{
			&mockTool{},
		})
		assert.NoError(t, err)

		callbacks := agent.Callbacks()

		assert.Equal(t, agent.opts.CallbackOptions.Callbacks, callbacks)
	})
}
