package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Anthropic satisfies the ChatModel interface.
var _ schema.ChatModel = (*Anthropic)(nil)

type AnthropicOptions struct {
	*schema.CallbackOptions
	Tokenizer schema.Tokenizer
	// Model name to use.
	ModelName string
	// Denotes the number of tokens to predict per generation.
	MaxTokens int
}

type Anthropic struct {
	schema.Tokenizer
	client *anthropic.Client
	opts   AnthropicOptions
}

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

func (cm *Anthropic) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	res, err := cm.client.Complete(ctx, &anthropic.CompletionRequest{
		Model:     cm.opts.ModelName,
		MaxTokens: cm.opts.MaxTokens,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: [][]schema.Generation{{newChatGeneraton(res.Completion)}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (cm *Anthropic) Type() string {
	return "chatmodel.Anthropic"
}

func (cm *Anthropic) Verbose() bool {
	return cm.opts.CallbackOptions.Verbose
}

func (cm *Anthropic) Callbacks() []schema.Callback {
	return cm.opts.CallbackOptions.Callbacks
}

func (cm *Anthropic) InvocationParams() map[string]any {
	return nil
}
