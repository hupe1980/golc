package model

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

func GeneratePrompt(ctx context.Context, model schema.Model, promptValues []schema.PromptValue, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
	if llm, ok := model.(schema.LLM); ok {
		prompts := util.Map(promptValues, func(value schema.PromptValue, _ int) string {
			return value.String()
		})

		return LLMGenerate(ctx, llm, prompts, optFns...)
	}

	if cm, ok := model.(schema.ChatModel); ok {
		messages := util.Map(promptValues, func(value schema.PromptValue, _ int) []schema.ChatMessage {
			return value.Messages()
		})

		return ChatModelGenerate(ctx, cm, messages, optFns...)
	}

	// TODO
	panic("invalid model type")
}

func LLMGenerate(ctx context.Context, model schema.LLM, prompts []string, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
	opts := schema.GenerateOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, model.Verbose())

	if err := cm.OnLLMStart(model.Type(), prompts); err != nil {
		return nil, err
	}

	result, err := model.Generate(ctx, prompts, opts.Stop)
	if err != nil {
		if cbErr := cm.OnLLMError(err); cbErr != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if err := cm.OnLLMEnd(result); err != nil {
		return nil, err
	}

	return result, nil
}

func ChatModelGenerate(ctx context.Context, cm schema.ChatModel, messages [][]schema.ChatMessage, optFns ...func(o *schema.GenerateOptions)) (*schema.LLMResult, error) {
	generations := [][]*schema.Generation{}

	for _, m := range messages {
		res, err := cm.Generate(ctx, m)
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &schema.LLMResult{
		Generations: generations,
	}, nil
}
