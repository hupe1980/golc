package llm

import (
	"context"

	"github.com/cohere-ai/cohere-go"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Cohere satisfies the LLM interface.
var _ schema.LLM = (*Cohere)(nil)

type CohereOptions struct {
	*schema.CallbackOptions
	Tokenizer  schema.Tokenizer
	Model      string
	Temperatur float32
}

type Cohere struct {
	schema.Tokenizer
	client *cohere.Client
	opts   CohereOptions
}

func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		Model: "medium",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewCohere("coheretext-50k")
		if tErr != nil {
			return nil, tErr
		}
	}

	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	return &Cohere{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

func (l *Cohere) Generate(ctx context.Context, prompts []string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	res, err := l.client.Generate(cohere.GenerateOptions{
		Model:         l.opts.Model,
		Prompt:        prompts[0],
		StopSequences: opts.Stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: [][]schema.Generation{{schema.Generation{Text: res.Generations[0].Text}}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (l *Cohere) Type() string {
	return "llm.Cohere"
}

func (l *Cohere) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

func (l *Cohere) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

func (l *Cohere) InvocationParams() map[string]any {
	return nil
}
