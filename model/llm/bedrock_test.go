package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/assert"
)

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
		bedrockModel, err := NewBedrock(client)
		assert.NoError(t, err)
		assert.Equal(t, "llm.Bedrock", bedrockModel.Type())
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
