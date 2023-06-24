package model

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type Options struct {
	Stop        []string
	Callbacks   []schema.Callback
	ParentRunID string
}

func GeneratePrompt(ctx context.Context, model schema.Model, promptValues []schema.PromptValue, optFns ...func(o *Options)) (*schema.LLMResult, error) {
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

func LLMGenerate(ctx context.Context, model schema.LLM, prompts []string, optFns ...func(o *Options)) (*schema.LLMResult, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, model.Callbacks(), model.Verbose(), func(mo *callback.ManagerOptions) {
		if opts.ParentRunID != "" {
			mo.ParentRunID = opts.ParentRunID
		}
	})

	rm, err := cm.OnLLMStart(model.Type(), prompts)
	if err != nil {
		return nil, err
	}

	result, err := model.Generate(ctx, prompts, func(o *schema.GenerateOptions) {
		o.CallbackManger = rm
		o.Stop = opts.Stop
	})
	if err != nil {
		if cbErr := rm.OnLLMError(err); cbErr != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if err := rm.OnLLMEnd(*result); err != nil {
		return nil, err
	}

	return result, nil
}

func ChatModelGenerate(ctx context.Context, model schema.ChatModel, messages [][]schema.ChatMessage, optFns ...func(o *Options)) (*schema.LLMResult, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	//cm := callback.NewManager(opts.Callbacks, model.Callbacks(), model.Verbose())

	generations := [][]schema.Generation{}

	for _, m := range messages {
		res, err := model.Generate(ctx, m)
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &schema.LLMResult{
		Generations: generations,
	}, nil
}
