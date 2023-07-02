package model

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

type Options struct {
	Stop        []string
	Callbacks   []schema.Callback
	ParentRunID string
	Functions   []schema.FunctionDefinition
}

func GeneratePrompt(ctx context.Context, model schema.Model, promptValues []schema.PromptValue, optFns ...func(o *Options)) (*schema.ModelResult, error) {
	if llm, ok := model.(schema.LLM); ok {
		prompts := make([]string, len(promptValues))
		for i, v := range promptValues {
			prompts[i] = v.String()
		}

		return LLMGenerate(ctx, llm, prompts, optFns...)
	}

	if cm, ok := model.(schema.ChatModel); ok {
		messages := make([]schema.ChatMessages, len(promptValues))
		for i, v := range promptValues {
			messages[i] = v.Messages()
		}

		return ChatModelGenerate(ctx, cm, messages, optFns...)
	}

	// TODO
	panic("invalid model type")
}

func LLMGenerate(ctx context.Context, model schema.LLM, prompts []string, optFns ...func(o *Options)) (*schema.ModelResult, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, model.Callbacks(), model.Verbose(), func(mo *callback.ManagerOptions) {
		if opts.ParentRunID != "" {
			mo.ParentRunID = opts.ParentRunID
		}
	})

	rm, err := cm.OnLLMStart(ctx, &schema.LLMStartManagerInput{
		LLMType:          model.Type(),
		Prompts:          prompts,
		InvocationParams: model.InvocationParams(),
	})
	if err != nil {
		return nil, err
	}

	result, err := model.Generate(ctx, prompts, func(o *schema.GenerateOptions) {
		o.CallbackManger = rm
		o.Stop = opts.Stop
	})
	if err != nil {
		if cbErr := rm.OnModelError(ctx, &schema.ModelErrorManagerInput{
			Error: err,
		}); cbErr != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if err := rm.OnModelEnd(ctx, &schema.ModelEndManagerInput{
		Result: result,
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func ChatModelGenerate(ctx context.Context, model schema.ChatModel, messages []schema.ChatMessages, optFns ...func(o *Options)) (*schema.ModelResult, error) {
	opts := Options{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, model.Callbacks(), model.Verbose(), func(mo *callback.ManagerOptions) {
		if opts.ParentRunID != "" {
			mo.ParentRunID = opts.ParentRunID
		}
	})

	rm, err := cm.OnChatModelStart(ctx, &schema.ChatModelStartManagerInput{
		ChatModelType:    model.Type(),
		Messages:         messages,
		InvocationParams: model.InvocationParams(),
	})
	if err != nil {
		return nil, err
	}

	generations := [][]schema.Generation{}

	for _, m := range messages {
		res, err := model.Generate(ctx, m, func(o *schema.GenerateOptions) {
			o.CallbackManger = rm
			o.Stop = opts.Stop
			o.Functions = opts.Functions
		})
		if err != nil {
			return nil, err
		}

		generations = append(generations, res.Generations...)
	}

	return &schema.ModelResult{
		Generations: generations,
	}, nil
}
