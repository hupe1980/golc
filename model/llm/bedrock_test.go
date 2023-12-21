package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/assert"
)

func TestBedrockInputOutputAdapter(t *testing.T) {
	t.Run("PrepareInput", func(t *testing.T) {
		tests := []struct {
			name         string
			provider     string
			prompt       string
			modelParams  map[string]interface{}
			expectedBody string
			expectedErr  string
		}{
			{
				name:     "PrepareInput for ai21",
				provider: "ai21",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","prompt":"Test prompt"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for amazon",
				provider: "amazon",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: `{"inputText":"Test prompt","textGenerationConfig":{"param1":"value1"}}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for anthropic",
				provider: "anthropic",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","max_tokens_to_sample":256,"prompt":"\n\nHuman:Test prompt\n\nAssistant:"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for cohere",
				provider: "cohere",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","prompt":"Test prompt"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for meta",
				provider: "meta",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: `{"param1":"value1","prompt":"Test prompt"}`,
				expectedErr:  "",
			},
			{
				name:     "PrepareInput for unsupported provider",
				provider: "xxx",
				prompt:   "Test prompt",
				modelParams: map[string]interface{}{
					"param1": "value1",
				},
				expectedBody: "",
				expectedErr:  "unsupported provider: xxx",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bioa := NewBedrockInputOutputAdapter(tt.provider)
				body, err := bioa.PrepareInput(tt.prompt, tt.modelParams)

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
				name:         "PrepareOutput for ai21",
				provider:     "ai21",
				response:     []byte(`{"completions":[{"data":{"text":"Generated text"}}]}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareOutput for amazon",
				provider:     "amazon",
				response:     []byte(`{"results":[{"OutputText":"Generated text"}]}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareOutput for anthropic",
				provider:     "anthropic",
				response:     []byte(`{"completion":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareOutput for cohere",
				provider:     "cohere",
				response:     []byte(`{"generations":[{"text":"Generated text"}]}`),
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
				name:         "PrepareStreamOutput for amazon",
				provider:     "amazon",
				response:     []byte(`{"outputText":"Streamed text"}`),
				expectedText: "Streamed text",
				expectedErr:  "",
			},
			{
				name:         "PrepareStreamOutput for anthropic",
				provider:     "anthropic",
				response:     []byte(`{"completion":"Generated text"}`),
				expectedText: "Generated text",
				expectedErr:  "",
			},
			{
				name:         "PrepareStreamOutput for cohere",
				provider:     "cohere",
				response:     []byte(`{"text":"Generated text"}`),
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

	t.Run("AI21", func(t *testing.T) {
		model, err := NewBedrockAI21(client)
		assert.NoError(t, err)

		t.Run("InvokeModel", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					b, err := json.Marshal(&ai21Output{
						Completions: []struct {
							Data struct {
								Text string `json:"text"`
							} `json:"data"`
						}{
							{
								Data: struct {
									Text string `json:"text"`
								}{
									Text: "Hello, how can I help you?",
								},
							},
						},
					})
					assert.NoError(t, err)

					return &bedrockruntime.InvokeModelOutput{
						Body: b,
					}, nil
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Amazon", func(t *testing.T) {
		model, err := NewBedrockAmazon(client)
		assert.NoError(t, err)

		t.Run("InvokeModel", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					b, err := json.Marshal(&amazonOutput{
						Results: []struct {
							OutputText       string `json:"outputText"`
							TokenCount       int    `json:"tokenCount"`
							CompletionReason string `json:"completionReason"`
						}{
							{
								OutputText: "Hello, how can I help you?",
							},
						},
					})
					assert.NoError(t, err)

					return &bedrockruntime.InvokeModelOutput{
						Body: b,
					}, nil
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Cohere", func(t *testing.T) {
		model, err := NewBedrockCohere(client)
		assert.NoError(t, err)

		t.Run("InvokeModel", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					b, err := json.Marshal(&cohereOutput{
						Generations: []struct {
							Text string `json:"text"`
						}{
							{
								Text: "Hello, how can I help you?",
							},
						},
					})
					assert.NoError(t, err)

					return &bedrockruntime.InvokeModelOutput{
						Body: b,
					}, nil
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Antrophic", func(t *testing.T) {
		model, err := NewBedrockAntrophic(client)
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

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				result, err := model.Generate(context.Background(), "Can you help me?")
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Meta", func(t *testing.T) {
		BedrockModel, err := NewBedrockMeta(client)
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

				result, err := BedrockModel.Generate(context.Background(), "Can you help me?")
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, result, "Expected non-nil result")
				assert.Len(t, result.Generations, 1, "Expected 1 generation")
				assert.Equal(t, "Hello, how can I help you?", result.Generations[0].Text, "Generated text does not match")
			})

			t.Run("Bedrock API error", func(t *testing.T) {
				client.createInvokeModelFn = func(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
					return nil, fmt.Errorf("bedrock api error")
				}

				result, err := BedrockModel.Generate(context.Background(), "Can you help me?")
				assert.Error(t, err, "Expected an error")
				assert.Nil(t, result, "Expected nil result")
			})
		})
	})

	t.Run("Type", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "amazon.titan-text-lite-v1")
		assert.NoError(t, err)
		assert.Equal(t, "llm.Bedrock", bedrockModel.Type())
	})

	t.Run("Callbacks", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "amazon.titan-text-lite-v1")
		assert.NoError(t, err)
		assert.Equal(t, bedrockModel.opts.CallbackOptions.Callbacks, bedrockModel.Callbacks())
	})

	t.Run("InvocationParams", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "amazon.titan-text-lite-v1", func(o *BedrockOptions) {
			o.ModelParams = map[string]any{
				"temperature": 0.7,
			}
		})
		assert.NoError(t, err)

		params := bedrockModel.InvocationParams()

		assert.Equal(t, "amazon.titan-text-lite-v1", params["model_id"])
		assert.Equal(t, 0.7, (params["model_params"].(map[string]any))["temperature"])
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
