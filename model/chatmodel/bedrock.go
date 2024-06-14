package chatmodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeDocument "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Bedrock satisfies the ChatModel interface.
var _ schema.ChatModel = (*Bedrock)(nil)

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	ConverseStream(ctx context.Context, params *bedrockruntime.ConverseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseStreamOutput, error)
	Converse(ctx context.Context, params *bedrockruntime.ConverseInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error)
}

type BedrockConverseOptions struct {
	Messages                          []bedrockruntimeTypes.Message
	ModelID                           *string
	AdditionalModelRequestFields      bedrockruntimeDocument.Interface
	AdditionalModelResponseFieldPaths []string
	InferenceConfig                   *bedrockruntimeTypes.InferenceConfiguration
	System                            []bedrockruntimeTypes.SystemContentBlock
	ToolConfig                        *bedrockruntimeTypes.ToolConfiguration
}

// BedrockAnthropicOptions contains options for configuring the Bedrock model with the "anthropic" provider.
type BedrockAnthropicOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// MaxTokensToSmaple sets the maximum number of tokens in the generated text.
	MaxTokensToSample int `map:"max_tokens_to_sample"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float32 `map:"temperature"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p,omitempty"`

	// TopK determines how the model selects tokens for output.
	TopK int `map:"top_k"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// NewBedrockAntrophic creates a new instance of Bedrock for the "anthropic" provider.
func NewBedrockAntrophic(client BedrockRuntimeClient, optFns ...func(o *BedrockAnthropicOptions)) (*Bedrock, error) {
	opts := BedrockAnthropicOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:           "anthropic.claude-v2", // https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature:       0.5,
		MaxTokensToSample: 256,
		TopP:              1,
		TopK:              250,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewClaude()
		if tErr != nil {
			return nil, tErr
		}
	}

	return NewBedrock(client, opts.ModelID, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.MaxTokens = aws.Int32(int32(opts.MaxTokensToSample))
		o.Temperature = aws.Float32(opts.Temperature)
		o.TopP = aws.Float32(opts.TopP)
		o.ModelParams = map[string]any{
			"top_k": opts.TopK,
		}
		o.Stream = opts.Stream
	})
}

// BedrockMetaOptions contains options for configuring the Bedrock model with the "meta" provider.
type BedrockMetaOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float32 `map:"temperature"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p,omitempty"`

	// MaxGenLen specify the maximum number of tokens to use in the generated response.
	MaxGenLen int `map:"max_gen_len"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// NewBedrockMeta creates a new instance of Bedrock for the "meta" provider.
func NewBedrockMeta(client BedrockRuntimeClient, optFns ...func(o *BedrockMetaOptions)) (*Bedrock, error) {
	opts := BedrockMetaOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelID:     "meta.llama2-70b-chat-v1", // https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
		Temperature: 0.5,
		TopP:        0.9,
		MaxGenLen:   512,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewGPT2()
		if tErr != nil {
			return nil, tErr
		}
	}

	return NewBedrock(client, opts.ModelID, func(o *BedrockOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.Tokenizer = opts.Tokenizer
		o.Temperature = aws.Float32(opts.Temperature)
		o.TopP = aws.Float32(opts.TopP)
		o.MaxTokens = aws.Int32(int32(opts.MaxGenLen))
		o.Stream = opts.Stream
	})
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens *int32

	// Stop is a list of sequences to stop the generation at.
	StopSequences []string

	// Temperature
	Temperature *float32

	// TopP
	TopP *float32

	// Additional model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// Bedrock is a model implementation of the schema.ChatModel interface for the Bedrock model.
type Bedrock struct {
	schema.Tokenizer
	client  BedrockRuntimeClient
	modelID string
	opts    BedrockOptions
}

// NewBedrock creates an instance of the Bedrock model.
func NewBedrock(client BedrockRuntimeClient, modelID string, optFns ...func(o *BedrockOptions)) (*Bedrock, error) {
	opts := BedrockOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelParams: make(map[string]any),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewGPT2()
		if tErr != nil {
			return nil, tErr
		}
	}

	return &Bedrock{
		Tokenizer: opts.Tokenizer,
		client:    client,
		modelID:   modelID,
		opts:      opts,
	}, nil
}

func (cm *Bedrock) PrepareInput(msgs schema.ChatMessages, params map[string]any) (*bedrockruntime.ConverseInput, error) {
	messages := make([]bedrockruntimeTypes.Message, 0, len(msgs))
	system := make([]bedrockruntimeTypes.SystemContentBlock, 0)

	for _, msg := range msgs {
		switch msg.Type() {
		case schema.ChatMessageTypeSystem:
			system = append(system, &bedrockruntimeTypes.SystemContentBlockMemberText{
				Value: msg.Content(),
			})
		case schema.ChatMessageTypeAI:
			messages = append(messages, bedrockruntimeTypes.Message{
				Role: bedrockruntimeTypes.ConversationRoleAssistant,
				Content: []bedrockruntimeTypes.ContentBlock{
					&bedrockruntimeTypes.ContentBlockMemberText{
						Value: msg.Content(),
					},
				},
			})
		default:
			messages = append(messages, bedrockruntimeTypes.Message{
				Role: bedrockruntimeTypes.ConversationRoleUser,
				Content: []bedrockruntimeTypes.ContentBlock{
					&bedrockruntimeTypes.ContentBlockMemberText{
						Value: msg.Content(),
					},
				},
			})
		}
	}

	var additionalModelRequestFields bedrockruntimeDocument.Interface

	if len(params) > 0 {
		additionalModelRequestFields = bedrockruntimeDocument.NewLazyDocument(params)
	}

	return &bedrockruntime.ConverseInput{
		Messages: messages,
		ModelId:  aws.String(cm.modelID),
		InferenceConfig: &bedrockruntimeTypes.InferenceConfiguration{
			MaxTokens:     cm.opts.MaxTokens,
			StopSequences: cm.opts.StopSequences,
			Temperature:   cm.opts.Temperature,
			TopP:          cm.opts.TopP,
		},
		System:                       system,
		AdditionalModelRequestFields: additionalModelRequestFields,
	}, nil
}

// Generate generates text based on the provided chat messages and options.
func (cm *Bedrock) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	params := util.CopyMap(cm.opts.ModelParams)

	input, err := cm.PrepareInput(messages, params)
	if err != nil {
		return nil, err
	}

	var completion string

	llmOutput := make(map[string]any)

	if cm.opts.Stream {
		input := &bedrockruntime.ConverseStreamInput{
			Messages:                     input.Messages,
			ModelId:                      input.ModelId,
			AdditionalModelRequestFields: input.AdditionalModelRequestFields,
			InferenceConfig:              input.InferenceConfig,
			System:                       input.System,
		}

		res, err := cm.client.ConverseStream(
			ctx,
			input,
		)
		if err != nil {
			return nil, err
		}

		stream := res.GetStream()

		defer stream.Close()

		tokens := []string{}

		for event := range stream.Events() {
			switch v := event.(type) {
			case *bedrockruntimeTypes.ConverseStreamOutputMemberContentBlockDelta:
				delta := v.Value.Delta

				token, ok := delta.(*bedrockruntimeTypes.ContentBlockDeltaMemberText)
				if !ok {
					return nil, fmt.Errorf("unexpected content type returned from bedrock: %T", v)
				}

				if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
					Token: token.Value,
				}); err != nil {
					return nil, err
				}

				tokens = append(tokens, token.Value)
			case *bedrockruntimeTypes.ConverseStreamOutputMemberMetadata:
				if v.Value.Usage == nil {
					continue
				}

				usage := v.Value.Usage

				if _, ok := llmOutput["input_tokens"]; !ok {
					llmOutput["input_tokens"] = *usage.InputTokens
				} else {
					llmOutput["input_tokens"] = llmOutput["input_tokens"].(int32) + *usage.InputTokens
				}

				if _, ok := llmOutput["output_tokens"]; !ok {
					llmOutput["output_tokens"] = *usage.OutputTokens
				} else {
					llmOutput["output_tokens"] = llmOutput["output_tokens"].(int32) + *usage.OutputTokens
				}

				if _, ok := llmOutput["tokens"]; !ok {
					llmOutput["tokens"] = *usage.TotalTokens
				} else {
					llmOutput["tokens"] = llmOutput["tokens"].(int32) + *usage.TotalTokens
				}
			}
		}

		completion = strings.Join(tokens, "")
	} else {
		res, err := cm.client.Converse(ctx, input)
		if err != nil {
			return nil, err
		}

		o, ok := res.Output.(*bedrockruntimeTypes.ConverseOutputMemberMessage)
		if !ok {
			return nil, fmt.Errorf("unexpected output type returned from bedrock: %T", res.Output)
		}

		var output string

		for _, block := range o.Value.Content {
			text, ok := block.(*bedrockruntimeTypes.ContentBlockMemberText)
			if !ok {
				return nil, fmt.Errorf("unexpected content type returned from bedrock: %T", block)
			}

			output += text.Value
		}

		completion = output

		if res.Usage != nil {
			llmOutput["input_tokens"] = *res.Usage.InputTokens
			llmOutput["output_tokens"] = *res.Usage.OutputTokens
			llmOutput["tokens"] = *res.Usage.TotalTokens
		}
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{newChatGeneraton(completion)},
		LLMOutput:   llmOutput,
	}, nil
}

// Type returns the type of the model.
func (cm *Bedrock) Type() string {
	return "chatmodel.Bedrock"
}

// Verbose returns the verbosity setting of the model.
func (cm *Bedrock) Verbose() bool {
	return cm.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Bedrock) Callbacks() []schema.Callback {
	return cm.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Bedrock) InvocationParams() map[string]any {
	params := util.StructToMap(cm.opts)
	params["model_id"] = cm.modelID

	return params
}
