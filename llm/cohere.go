package llm

import (
	"context"

	"github.com/cohere-ai/cohere-go"
	"github.com/hupe1980/golc"
)

// Compile time check to ensure Cohere satisfies the llm interface.
var _ golc.LLM = (*Cohere)(nil)

type CohereOptions struct {
	Model      string
	Temperatur float32
}

type Cohere struct {
	*LLM
	client *cohere.Client
	model  string
}

func NewCohere(apiKey string) (*Cohere, error) {
	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	cohere := &Cohere{
		client: client,
		model:  "medium",
	}

	cohere.LLM = NewLLM(cohere.generate)

	return cohere, nil
}

func (co *Cohere) generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	res, err := co.client.Generate(cohere.GenerateOptions{
		Model:  "medium",
		Prompt: prompts[0],
	})
	if err != nil {
		return nil, err
	}

	return &golc.LLMResult{
		Generations: [][]golc.Generation{{golc.Generation{Text: res.Generations[0].Text}}},
		LLMOutput:   map[string]any{},
	}, nil
}
