package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/tokenizer"
)

type AnthropicOptions struct {
	// Model name to use.
	ModelName string
	// Denotes the number of tokens to predict per generation.
	MaxTokens int
}

type Anthropic struct {
	*ChatModel
	golc.Tokenizer
	client *anthropic.Client
	opts   AnthropicOptions
}

func NewAnthropic(apiKey string) (*Anthropic, error) {
	opts := AnthropicOptions{
		ModelName: "claude-v1",
		MaxTokens: 256,
	}

	a := &Anthropic{
		Tokenizer: tokenizer.NewSimple(),
		client:    anthropic.New(apiKey),
		opts:      opts,
	}

	a.ChatModel = NewChatModel(a.generate)

	return a, nil
}

func (a *Anthropic) generate(ctx context.Context, messages []golc.ChatMessage) (*golc.LLMResult, error) {
	res, err := a.client.Complete(ctx, &anthropic.CompletionRequest{
		Model:     a.opts.ModelName,
		MaxTokens: a.opts.MaxTokens,
	})
	if err != nil {
		return nil, err
	}

	return &golc.LLMResult{
		Generations: [][]*golc.Generation{{newChatGeneraton(res.Completion)}},
		LLMOutput:   map[string]any{},
	}, nil
}
