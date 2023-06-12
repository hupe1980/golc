package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

type GenerateFunc func(ctx context.Context, prompts []string) (*golc.LLMResult, error)

type LLM struct {
	generateFunc GenerateFunc
}

func NewLLM(generateFunc GenerateFunc) *LLM {
	return &LLM{
		generateFunc: generateFunc,
	}
}

func (b *LLM) Generate(ctx context.Context, prompts []string) (*golc.LLMResult, error) {
	return b.generateFunc(ctx, prompts)
}

func (b *LLM) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue) string {
		return value.String()
	})

	return b.generateFunc(ctx, prompts)
}

func (b *LLM) Call(ctx context.Context, prompt string) (string, error) {
	result, err := b.Generate(ctx, []string{prompt})
	if err != nil {
		return "", err
	}

	return result.Generations[0][0].Text, nil
}

func (b *LLM) Predict(ctx context.Context, text string) (string, error) {
	return b.Call(ctx, text)
}

func (b *LLM) PredictMessages(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return nil, err
	}

	prediction, err := b.Call(ctx, text)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(prediction), nil
}
