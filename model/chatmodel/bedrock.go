package chatmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
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
func (bioa *BedrockInputOutputAdapter) PrepareInput(messages schema.ChatMessages, modelParams map[string]any) ([]byte, error) {
	body := modelParams

	switch bioa.provider {
	case "anthropic":
		p, err := convertMessagesToAnthropicPrompt(messages)
		if err != nil {
			return nil, err
		}

		body["prompt"] = p

		if _, ok := body["max_tokens_to_sample"]; !ok {
			body["max_tokens_to_sample"] = 256
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", bioa.provider)
	}

	return json.Marshal(body)
}

// anthropicOutput is a struct representing the output structure for the "anthropic" provider.
type anthropicOutput struct {
	Completion string `json:"completion"`
}

// PrepareOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareOutput(response *bedrockruntime.InvokeModelOutput) (string, error) {
	switch bioa.provider {
	case "anthropic":
		output := &anthropicOutput{}
		if err := json.Unmarshal(response.Body, output); err != nil {
			return "", err
		}

		return output.Completion, nil
	}

	return "", fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`
}

// Bedrock is a model implementation of the schema.ChatModel interface for the Bedrock model.
type Bedrock struct {
	schema.Tokenizer
	client BedrockRuntimeClient
	opts   BedrockOptions
}

// NewBedrock creates an instance of the Bedrock model.
func NewBedrock(client BedrockRuntimeClient, optFns ...func(o *BedrockOptions)) (*Bedrock, error) {
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

	if len(opts.Stop) > 0 {
		params["stop_sequences"] = opts.Stop
	}

	bioa := NewBedrockInputOutputAdapter(cm.getProvider())

	body, err := bioa.PrepareInput(messages, params)
	if err != nil {
		return nil, err
	}

	res, err := cm.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(cm.opts.ModelID),
		Body:        body,
		Accept:      aws.String("application/json"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	completion, err := bioa.PrepareOutput(res)
	if err != nil {
		return nil, err
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
	return []schema.Callback{}
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Bedrock) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}

// getProvider returns the provider of the model based on the model ID.
func (cm *Bedrock) getProvider() string {
	return strings.Split(cm.opts.ModelID, ".")[0]
}
