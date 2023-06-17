package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/util"
)

type generateFunc func(ctx context.Context, prompts []string) (*golc.LLMResult, error)

type llm struct {
	llmName      string
	generateFunc generateFunc
	verbose      bool
}

func newLLM(llmName string, generateFunc generateFunc, verbose bool) *llm {
	return &llm{
		llmName:      llmName,
		generateFunc: generateFunc,
		verbose:      verbose,
	}
}

func (l *llm) Generate(ctx context.Context, prompts []string, optFns ...func(o *golc.GenerateOptions)) (*golc.LLMResult, error) {
	opts := golc.GenerateOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, l.verbose)

	if err := cm.OnLLMStart(l.llmName, prompts); err != nil {
		return nil, err
	}

	result, err := l.generateFunc(ctx, prompts)
	if err != nil {
		if cbErr := cm.OnLLMError(err); err != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if err := cm.OnLLMEnd(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (l *llm) GeneratePrompt(ctx context.Context, promptValues []golc.PromptValue, optFns ...func(o *golc.GenerateOptions)) (*golc.LLMResult, error) {
	prompts := util.Map(promptValues, func(value golc.PromptValue, _ int) string {
		return value.String()
	})

	return l.Generate(ctx, prompts, optFns...)
}

func (l *llm) Predict(ctx context.Context, text string, optFns ...func(o *golc.GenerateOptions)) (string, error) {
	result, err := l.Generate(ctx, []string{text}, optFns...)
	if err != nil {
		return "", err
	}

	return result.Generations[0][0].Text, nil
}

func (l *llm) PredictMessages(ctx context.Context, messages []golc.ChatMessage, optFns ...func(o *golc.GenerateOptions)) (golc.ChatMessage, error) {
	text, err := golc.StringifyChatMessages(messages)
	if err != nil {
		return nil, err
	}

	prediction, err := l.Predict(ctx, text, optFns...)
	if err != nil {
		return nil, err
	}

	return golc.NewAIChatMessage(prediction), nil
}

type callbackOptions struct {
	Callbacks []golc.Callback
	Verbose   bool
}
