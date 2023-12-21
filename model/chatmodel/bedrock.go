package chatmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Bedrock satisfies the ChatModel interface.
var _ schema.ChatModel = (*Bedrock)(nil)

// BedrockInputOutputAdapter is a helper struct for preparing input and handling output for Bedrock model.
type BedrockInputOutputAdapter struct {
	provider string
}

// NewBedrockInputOutputAdpter creates a new instance of BedrockInputOutputAdpter.
func NewBedrockInputOutputAdapter(provider string) *BedrockInputOutputAdapter {
	return &BedrockInputOutputAdapter{
		provider: provider,
	}
}

// PrepareInput prepares the input for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareInput(messages schema.ChatMessages, modelParams map[string]any, stop []string) ([]byte, error) {
	body := modelParams

	switch bioa.provider {
	case "anthropic":
		p, err := convertMessagesToAnthropicPrompt(messages)
		if err != nil {
			return nil, err
		}

		body["prompt"] = p

		if len(stop) > 0 {
			body["stop_sequences"] = stop
		}

		if _, ok := body["max_tokens_to_sample"]; !ok {
			body["max_tokens_to_sample"] = 256
		}
	case "meta":
		p, err := convertMessagesToMetaPrompt(messages)
		if err != nil {
			return nil, err
		}

		body["prompt"] = p
	default:
		return nil, fmt.Errorf("unsupported provider: %s", bioa.provider)
	}

	return json.Marshal(body)
}

// anthropicOutput is a struct representing the output structure for the "anthropic" provider.
type anthropicOutput struct {
	Completion string `json:"completion"`
}

// metaOutput is a struct representing the output structure for the "meta" provider.
type metaOutput struct {
	Generation string `json:"generation"`
}

// PrepareOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareOutput(response []byte) (string, error) {
	switch bioa.provider {
	case "anthropic":
		output := &anthropicOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Completion, nil
	case "meta":
		output := &metaOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Generation, nil
	}

	return "", fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// anthropicStreamOutput is a struct representing the stream output structure for the "anthropic" provider.
type anthropicStreamOutput struct {
	Completion string `json:"completion"`
}

// metaStreamOutput is a struct representing the stream output structure for the "meta" provider.
type metaStreamOutput struct {
	Generation string `json:"generation"`
}

// PrepareStreamOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareStreamOutput(response []byte) (string, error) {
	switch bioa.provider {
	case "anthropic":
		output := &anthropicStreamOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Completion, nil
	case "meta":
		output := &metaStreamOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return "", err
		}

		return output.Generation, nil
	}

	return "", fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
	InvokeModelWithResponseStream(ctx context.Context, params *bedrockruntime.InvokeModelWithResponseStreamInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error)
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
		ModelID:           "anthropic.claude-v2", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
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
		o.ModelParams = map[string]any{
			"max_tokens_to_sample": opts.MaxTokensToSample,
			"temperature":          opts.Temperature,
			"top_p":                opts.TopP,
			"top_k":                opts.TopK,
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
		ModelID:     "meta.llama2-70b-chat-v1", //https://docs.aws.amazon.com/bedrock/latest/userguide/model-ids-arns.html
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
		o.ModelParams = map[string]any{
			"temperature": opts.Temperature,
			"top_p":       opts.TopP,
			"max_gen_len": opts.MaxGenLen,
		}
		o.Stream = opts.Stream
	})
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model params to use.
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

// Generate generates text based on the provided chat messages and options.
func (cm *Bedrock) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	params := util.CopyMap(cm.opts.ModelParams)

	bioa := NewBedrockInputOutputAdapter(cm.getProvider())

	body, err := bioa.PrepareInput(messages, params, opts.Stop)
	if err != nil {
		return nil, err
	}

	var completion string

	if cm.opts.Stream {
		res, err := cm.client.InvokeModelWithResponseStream(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
			ModelId:     aws.String(cm.modelID),
			Body:        body,
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			return nil, err
		}

		stream := res.GetStream()

		defer stream.Close()

		tokens := []string{}

		for event := range stream.Events() {
			switch v := event.(type) {
			case *bedrockruntimeTypes.ResponseStreamMemberChunk:
				token, err := bioa.PrepareStreamOutput(v.Value.Bytes)
				if err != nil {
					return nil, err
				}

				if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
					Token: token,
				}); err != nil {
					return nil, err
				}

				tokens = append(tokens, token)
			}
		}

		completion = strings.Join(tokens, "")
	} else {
		res, err := cm.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(cm.modelID),
			Body:        body,
			Accept:      aws.String("application/json"),
			ContentType: aws.String("application/json"),
		})
		if err != nil {
			return nil, err
		}

		output, err := bioa.PrepareOutput(res.Body)
		if err != nil {
			return nil, err
		}

		completion = output
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{newChatGeneraton(completion)},
		LLMOutput:   map[string]any{},
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

// getProvider returns the provider of the model based on the model ID.
func (cm *Bedrock) getProvider() string {
	return strings.Split(cm.modelID, ".")[0]
}
