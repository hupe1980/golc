package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
)

type LLMChainOptions struct {
	OutputKey string
}

type LLMChain struct {
	*Chain
	llm    golc.LLM
	prompt *prompt.Template
	opts   LLMChainOptions
}

func NewLLMChain(llm golc.LLM, prompt *prompt.Template) (*LLMChain, error) {
	opts := LLMChainOptions{
		OutputKey: "text",
	}

	llmChain := &LLMChain{
		prompt: prompt,
		llm:    llm,
		opts:   opts,
	}

	llmChain.Chain = NewChain(llmChain.call)

	return llmChain, nil
}

func (c *LLMChain) Predict(ctx context.Context, values golc.ChainValues) (string, error) {
	output, err := c.Call(ctx, values)
	if err != nil {
		return "", err
	}

	return output[c.opts.OutputKey].(string), err
}

func (c *LLMChain) call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	promptValue, err := c.prompt.FormatPrompt(values)
	if err != nil {
		return nil, err
	}

	res, err := c.llm.GeneratePrompt(ctx, []golc.PromptValue{promptValue})
	if err != nil {
		return nil, err
	}

	output, err := c.getFinalOutput(res.Generations[0])
	if err != nil {
		return nil, err
	}

	return golc.ChainValues{
		c.opts.OutputKey: output,
	}, nil
}

func (c *LLMChain) getFinalOutput(generations []golc.Generation) (any, error) { // nolint unparam
	completion := generations[0].Text
	// TODO Outputparser
	return completion, nil
}
