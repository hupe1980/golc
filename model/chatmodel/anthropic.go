package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Anthropic satisfies the ChatModel interface.
var _ schema.ChatModel = (*Anthropic)(nil)

type AnthropicOptions struct {
	*schema.CallbackOptions
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

func NewAnthropic(apiKey string) (*Anthropic, error) {
	opts := AnthropicOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName: "claude-v1",
		MaxTokens: 256,
	}

	return &Anthropic{
		Tokenizer: tokenizer.NewSimple(),
		client:    anthropic.New(apiKey),
		opts:      opts,
	}, nil
}

func (cm *Anthropic) Generate(ctx context.Context, messages schema.ChatMessages) (*schema.LLMResult, error) {
	res, err := cm.client.Complete(ctx, &anthropic.CompletionRequest{
		Model:     cm.opts.ModelName,
		MaxTokens: cm.opts.MaxTokens,
	})
	if err != nil {
		return nil, err
	}

	return &schema.LLMResult{
		Generations: [][]*schema.Generation{{newChatGeneraton(res.Completion)}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (cm *Anthropic) Type() string {
	return "Anthropic"
}

func (cm *Anthropic) Verbose() bool {
	return cm.opts.CallbackOptions.Verbose
}

func (cm *Anthropic) Callbacks() []schema.Callback {
	return cm.opts.CallbackOptions.Callbacks
}
