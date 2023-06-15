package llm

import (
	"context"

	"github.com/cohere-ai/cohere-go"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Cohere satisfies the llm interface.
var _ golc.LLM = (*Cohere)(nil)

type CohereOptions struct {
	Model      string
	Temperatur float32
}

type Cohere struct {
	*tokenizer
	client *cohere.Client
	opts   CohereOptions
}

func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		Model: "medium",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	cohere := &Cohere{
		client: client,
		opts:   opts,
	}

	return cohere, nil
}

func (co *Cohere) Generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	res, err := co.client.Generate(cohere.GenerateOptions{
		Model:  co.opts.Model,
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

func (co *Cohere) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue, _ int) string {
		return value.String()
	})

	return co.Generate(ctx, prompts)
}

func (co *Cohere) Predict(ctx context.Context, text string) (string, error) {
	result, err := co.Generate(ctx, []string{text})
	if err != nil {
		return "", err
	}

	return result.Generations[0][0].Text, nil
}

func (co *Cohere) PredictMessages(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return nil, err
	}

	prediction, err := co.Predict(ctx, text)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(prediction), nil
}
