package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

type GenerateFunc func(ctx context.Context, messages []golc.ChatMessage, optFns ...func(o *golc.GenerateOptions)) (*golc.LLMResult, error)

type ChatModel struct {
	generateFunc GenerateFunc
}

func NewChatModel(generateFunc GenerateFunc) *ChatModel {
	return &ChatModel{
		generateFunc: generateFunc,
	}
}

func (b *ChatModel) Generate(ctx context.Context, messages [][]golc.ChatMessage, optFns ...func(o *golc.GenerateOptions)) (*golc.LLMResult, error) {
	generations := [][]*golc.Generation{}

	for _, m := range messages {
		res, err := b.generateFunc(ctx, m, optFns...)
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &golc.LLMResult{
		Generations: generations,
	}, nil
}

func (b *ChatModel) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue, optFns ...func(o *golc.GenerateOptions)) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue, _ int) []golc.ChatMessage {
		return value.Messages()
	})

	return b.Generate(ctx, prompts, optFns...)
}

func (b *ChatModel) Predict(ctx context.Context, text string, optFns ...func(o *golc.GenerateOptions)) (string, error) {
	message := golc.NewHumanChatMessage(text)

	result, err := b.PredictMessages(ctx, []golc.ChatMessage{message}, optFns...)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func (b *ChatModel) PredictMessages(ctx context.Context, messages []golc.ChatMessage, optFns ...func(o *golc.GenerateOptions)) (golc.ChatMessage, error) {
	result, err := b.Generate(ctx, [][]golc.ChatMessage{messages}, optFns...)
	if err != nil {
		return nil, err
	}

	return result.Generations[0][0].Message, nil
}

func newChatGeneraton(text string) *golc.Generation {
	return &golc.Generation{
		Text:    text,
		Message: golc.NewAIChatMessage(text),
	}
}
