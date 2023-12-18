package chatmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestBedrockInputOutputAdapter(t *testing.T) {
	t.Run("PrepareInput", func(t *testing.T) {
		tests := []struct {
			name         string
			provider     string
			messages     schema.ChatMessages
			modelParams  map[string]any
			expectedBody string
			expectedErr  string
		}{
			{
				name:     "PrepareInput for anthropic",
				provider: "anthropic",
				messages: schema.ChatMessages{schema.NewHumanChatMessage("Test prompt")},
				modelParams: map[string]any{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","max_tokens_to_sample":256,"prompt":"\n\nHuman: Test prompt\n\nAssistant:"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for meta",
				provider: "meta",
				messages: schema.ChatMessages{schema.NewHumanChatMessage("Test prompt")},
				modelParams: map[string]any{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","prompt":"[INST] Test prompt [/INST]"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for unsupported provider",
				provider: "xxx",
				messages: nil,
				modelParams: map[string]any{
					"param1": "value1",
				},
				expectedBody: "",
				expectedErr:  "unsupported provider: xxx",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bioa := NewBedrockInputOutputAdapter(tt.provider)
				body, err := bioa.PrepareInput(tt.messages, tt.modelParams, []string{})

				if tt.expectedErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.JSONEq(t, tt.expectedBody, string(body))
				}
			})
		}
	})

	t.Run("PrepareOutput", func(t *testing.T) {
		tests := []struct {
			name         string
			provider     string
			response     []byte
			expectedText string
			expectedErr  string
		}{
			{
				name:         "PrepareOutput for anthropic",
				provider:     "anthropic",
				response:     []byte(`{"completion":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareOutput for meta",
				provider:     "meta",
				response:     []byte(`{"generation":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareOutput for unsupported provider",
				provider:     "xxx",
				response:     nil,
				expectedText: "",
				expectedErr:  "unsupported provider: xxx",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bioa := NewBedrockInputOutputAdapter(tt.provider)
				text, err := bioa.PrepareOutput(tt.response)

				if tt.expectedErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedText, text)
				}
			})
		}
	})

	t.Run("PrepareStreamOutput", func(t *testing.T) {
		tests := []struct {
			name         string
			provider     string
			response     []byte
			expectedText string
			expectedErr  string
		}{
			{
				name:         "PrepareStreamOutput for anthropic",
				provider:     "anthropic",
				response:     []byte(`{"completion":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareStreamOutput for meta",
				provider:     "meta",
				response:     []byte(`{"generation":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareStreamOutput for unsupported provider",
				provider:     "xxx",
				response:     nil,
				expectedText: "",
				expectedErr:  "unsupported provider: xxx",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bioa := NewBedrockInputOutputAdapter(tt.provider)
				text, err := bioa.PrepareStreamOutput(tt.response)

				if tt.expectedErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedText, text)
				}
			})
		}
	})
}

func TestBedrock(t *testing.T) {
	client := &mockBedrockClient{}

	t.Run("Antrophic", func(t *testing.T) {
		bedrockModel, err := NewBedrockAntrophic(client)
		assert.NoError(t, err)

		t.Run("InvokeModel", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					b, err := json.Marshal(&anthropicOutput{
						Completion: "Hello, how can I help you?",
					})
					assert.NoError(t, err)

					return &bedrockruntime.InvokeModelOutput{
						Body: b,
					}, nil
				}

				// Define chat messages
				chatMessages := []schema.ChatMessage{
					schema.NewAIChatMessage("Hi"),
					schema.NewHumanChatMessage("Can you help me?"),
				}

				result, err := bedrockModel.Generate(context.Background(), chatMessages)
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				// Define chat messages
				chatMessages := []schema.ChatMessage{
					schema.NewAIChatMessage("Hi"),
					schema.NewHumanChatMessage("Can you help me?"),
				}

				// Generate text
				result, err := bedrockModel.Generate(context.Background(), chatMessages)
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Meta", func(t *testing.T) {
		model, err := NewBedrockMeta(client)
		assert.NoError(t, err)

		t.Run("InvokeModel", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					b, err := json.Marshal(&metaOutput{
						Generation: "Hello, how can I help you?",
					})
					assert.NoError(t, err)

					return &bedrockruntime.InvokeModelOutput{
						Body: b,
					}, nil
				}

				// Define chat messages
				chatMessages := []schema.ChatMessage{
					schema.NewAIChatMessage("Hi"),
					schema.NewHumanChatMessage("Can you help me?"),
				}

				result, err := model.Generate(context.Background(), chatMessages)
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				// Define chat messages
				chatMessages := []schema.ChatMessage{
					schema.NewAIChatMessage("Hi"),
					schema.NewHumanChatMessage("Can you help me?"),
				}

				// Generate text
				result, err := model.Generate(context.Background(), chatMessages)
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Type", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client)
		assert.NoError(t, err)
		assert.Equal(t, "chatmodel.Bedrock", bedrockModel.Type())
	})

	t.Run("Callbacks", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client)
		assert.NoError(t, err)
		assert.Equal(t, bedrockModel.opts.CallbackOptions.Callbacks, bedrockModel.Callbacks())
	})

	t.Run("InvocationParams", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, func(o *BedrockOptions) {
			o.ModelID = "foo.bar"
		})
		assert.NoError(t, err)

		params := bedrockModel.InvocationParams()

		assert.Equal(t, "foo.bar", params["model_id"])
	})
}

// mockBedrockClient is a mock implementation of the BedrockClient interface for testing.
type mockBedrockClient struct {
	createInvokeModelFn func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

func (m *mockBedrockClient) InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
	return m.createInvokeModelFn(ctx, params)
}

func (m *mockBedrockClient) InvokeModelWithResponseStream(ctx context.Context, params *bedrockruntime.InvokeModelWithResponseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	return nil, nil
}
