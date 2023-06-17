package chatmodel

import "github.com/hupe1980/golc/integration/anthropic"

type AnthropicOptions struct{}

type Anthropic struct {
	client *anthropic.Client
	opts   AnthropicOptions
}

func NewAnthropic(apiKey string) (*Anthropic, error) {
	opts := AnthropicOptions{}

	return &Anthropic{
		client: anthropic.New(apiKey),
		opts:   opts,
	}, nil
}
