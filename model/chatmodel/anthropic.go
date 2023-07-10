package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Anthropic satisfies the ChatModel interface.
var _ schema.ChatModel = (*Anthropic)(nil)

// AnthropicOptions contains options for configuring the Anthropic chat model.
type AnthropicOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
	// Model name to use.
	ModelName string `map:"model_name,omitempty"`
	// Temperature parameter controls the randomness of the generation output.
	Temperature float64 `map:"temperature,omitempty"`
	// Denotes the number of tokens to predict per generation.
	MaxTokens int `map:"max_tokens,omitempty"`
	// TopK parameter specifies the number of highest probability tokens to consider for generation.
	TopK int `map:"top_k,omitempty"`
	// TopP parameter specifies the cumulative probability threshold for generating tokens.
	TopP float64 `map:"top_p,omitempty"`
}

// Anthropic is a chat model based on the Anthropic API.
type Anthropic struct {
	schema.Tokenizer
	client *anthropic.Client
	opts   AnthropicOptions
}

// NewAnthropic creates a new instance of the Anthropic chat model with the provided options.
func NewAnthropic(apiKey string, optFns ...func(o *AnthropicOptions)) (*Anthropic, error) {
	opts := AnthropicOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName: "claude-v1",
		MaxTokens: 256,
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

	return &Anthropic{
		Tokenizer: opts.Tokenizer,
		client:    anthropic.New(apiKey),
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided chat messages and options.
func (cm *Anthropic) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	prompt, err := messages.Format()
	if err != nil {
		return nil, err
	}

	res, err := cm.client.Complete(ctx, &anthropic.CompletionRequest{
		Prompt:      prompt,
		Model:       cm.opts.ModelName,
		Temperature: cm.opts.Temperature,
		MaxTokens:   cm.opts.MaxTokens,
		TopK:        cm.opts.TopK,
		TopP:        cm.opts.TopP,
		Stop:        opts.Stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{newChatGeneraton(res.Completion)},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (cm *Anthropic) Type() string {
	return "chatmodel.Anthropic"
}

// Verbose returns the verbosity setting of the model.
func (cm *Anthropic) Verbose() bool {
	return cm.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Anthropic) Callbacks() []schema.Callback {
	return cm.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Anthropic) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}
