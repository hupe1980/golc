package chatmodel

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestBedrock(t *testing.T) {
	client := &mockBedrockClient{}

	t.Run("Antrophic", func(t *testing.T) {
		bedrockModel, err := NewBedrockAntrophic(client)
		assert.NoError(t, err)

		t.Run("Converse", func(t *testing.T) {
			t.Run("Successful generation", func(t *testing.T) {
				client.createConverseFn = func(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
					messages := []bedrockruntimeTypes.ContentBlock{
						&bedrockruntimeTypes.ContentBlockMemberText{
							Value: "Hello, how can I help you?",
						},
					}

					return &bedrockruntime.ConverseOutput{
						Output: &bedrockruntimeTypes.ConverseOutputMemberMessage{
							Value: bedrockruntimeTypes.Message{
								Content: messages,
							},
						},
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
				client.createConverseFn = func(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
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
				client.createConverseFn = func(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
					messages := []bedrockruntimeTypes.ContentBlock{
						&bedrockruntimeTypes.ContentBlockMemberText{
							Value: "Hello, how can I help you?",
						},
					}

					return &bedrockruntime.ConverseOutput{
						Output: &bedrockruntimeTypes.ConverseOutputMemberMessage{
							Value: bedrockruntimeTypes.Message{
								Content: messages,
							},
						},
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
				client.createConverseFn = func(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
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
		bedrockModel, err := NewBedrock(client, "anthropic.claude-v2")
		assert.NoError(t, err)
		assert.Equal(t, "chatmodel.Bedrock", bedrockModel.Type())
	})

	t.Run("Callbacks", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "anthropic.claude-v2")
		assert.NoError(t, err)
		assert.Equal(t, bedrockModel.opts.CallbackOptions.Callbacks, bedrockModel.Callbacks())
	})

	t.Run("Verbose", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "anthropic.claude-v2")
		assert.NoError(t, err)
		assert.Equal(t, bedrockModel.opts.CallbackOptions.Verbose, bedrockModel.Verbose())
	})

	t.Run("InvocationParams", func(t *testing.T) {
		bedrockModel, err := NewBedrock(client, "anthropic.claude-v2", func(o *BedrockOptions) {
			o.ModelParams = map[string]any{
				"temperature": 0.7,
			}
		})
		assert.NoError(t, err)

		params := bedrockModel.InvocationParams()

		assert.Equal(t, "anthropic.claude-v2", params["model_id"])
		assert.Equal(t, 0.7, (params["model_params"].(map[string]any))["temperature"])
	})
}

// mockBedrockClient is a mock implementation of the BedrockClient interface for testing.
type mockBedrockClient struct {
	createConverseFn func(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error)
}

func (m *mockBedrockClient) Converse(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error) {
	return m.createConverseFn(ctx, params)
}

func (m *mockBedrockClient) ConverseStream(ctx context.Context, params *bedrockruntime.ConverseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseStreamOutput, error) {
	return nil, nil
}
