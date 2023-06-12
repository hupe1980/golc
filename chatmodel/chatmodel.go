package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

type GenerateFunc func(ctx context.Context, messages []golc.ChatMessage) (*golc.LLMResult, error)

type ChatModel struct {
	generateFunc GenerateFunc
}

func NewChatModel(generateFunc GenerateFunc) *ChatModel {
	return &ChatModel{
		generateFunc: generateFunc,
	}
}

func (b *ChatModel) Generate(ctx context.Context, messages [][]golc.ChatMessage) (*golc.LLMResult, error) {
	generations := [][]golc.Generation{}

	for _, m := range messages {
		res, err := b.generateFunc(ctx, m)
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &golc.LLMResult{
		Generations: generations,
	}, nil
}

func (b *ChatModel) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue) []golc.ChatMessage {
		return value.Messages()
	})

	return b.Generate(ctx, prompts)
}

func (b *ChatModel) Call(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	result, err := b.Generate(ctx, [][]golc.ChatMessage{messages})
	if err != nil {
		return nil, err
	}

	return result.Generations[0][0].Message, nil
}

func (b *ChatModel) CallPrompt(ctx context.Context, promptValue golc.PromptValue) (golc.ChatMessage, error) {
	promptMessages := promptValue.Messages()
	return b.Call(ctx, promptMessages)
}

func (b *ChatModel) Predict(ctx context.Context, text string) (string, error) {
	message := golc.NewHumanChatMessage(text)

	result, err := b.Call(ctx, []golc.ChatMessage{message})
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func (b *ChatModel) PredictMessages(ctx context.Context, messages []golc.ChatMessage) (golc.ChatMessage, error) {
	return b.Call(ctx, messages)
}
