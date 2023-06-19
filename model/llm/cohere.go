package llm

import (
	"context"

	"github.com/cohere-ai/cohere-go"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Cohere satisfies the LLM interface.
var _ schema.LLM = (*Cohere)(nil)

type CohereOptions struct {
	Model      string
	Temperatur float32
	callbackOptions
}

type Cohere struct {
	*llm
	schema.Tokenizer
	client *cohere.Client
	opts   CohereOptions
}

func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		Model: "medium",
		callbackOptions: callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	cohere := &Cohere{
		Tokenizer: tokenizer.NewSimple(),
		client:    client,
		opts:      opts,
	}

	cohere.llm = newLLM("Cohere", cohere.generate, opts.Verbose)

	return cohere, nil
}

func (co *Cohere) generate(ctx context.Context, prompts []string, stop []string) (*schema.LLMResult, error) {
	res, err := co.client.Generate(cohere.GenerateOptions{
		Model:         co.opts.Model,
		Prompt:        prompts[0],
		StopSequences: stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.LLMResult{
		Generations: [][]*schema.Generation{{&schema.Generation{Text: res.Generations[0].Text}}},
		LLMOutput:   map[string]any{},
	}, nil
}