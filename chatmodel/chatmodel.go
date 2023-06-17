package chatmodel

import (
	"context"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type GenerateFunc func(ctx context.Context, messages []schema.ChatMessage, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error)

type ChatModel struct {
	generateFunc GenerateFunc
}

func NewChatModel(generateFunc GenerateFunc) *ChatModel {
	return &ChatModel{
		generateFunc: generateFunc,
	}
}

func (b *ChatModel) Generate(ctx context.Context, messages [][]schema.ChatMessage, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
	generations := [][]*schema.Generation{}

	for _, m := range messages {
		res, err := b.generateFunc(ctx, m, optFns...)
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &schema.LLMResult{
		Generations: generations,
	}, nil
}

func (b *ChatModel) GeneratePrompt(ctx context.Context, promptValues []schema.PromptValue, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
	prompts := util.Map(promptValues, func(value schema.PromptValue, _ int) []schema.ChatMessage {
		return value.Messages()
	})

	return b.Generate(ctx, prompts, optFns...)
}

func (b *ChatModel) Predict(ctx context.Context, text string, optFns ...func(o *schema.GenerateOptions)) (string, error) {
	message := schema.NewHumanChatMessage(text)

	result, err := b.PredictMessages(ctx, []schema.ChatMessage{message}, optFns...)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func (b *ChatModel) PredictMessages(ctx context.Context, messages []schema.ChatMessage, optFns ...func(o *schema.GenerateOptions)) (schema.ChatMessage, error) {
	result, err := b.Generate(ctx, [][]schema.ChatMessage{messages}, optFns...)
	if err != nil {
		return nil, err
	}

	return result.Generations[0][0].Message, nil
}

func newChatGeneraton(text string) *schema.Generation {
	return &schema.Generation{
		Text:    text,
		Message: schema.NewAIChatMessage(text),
	}
}
