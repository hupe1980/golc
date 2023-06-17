package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type LLMChainOptions struct {
	callbackOptions
	OutputKey    string
	OutputParser schema.OutputParser[any]
}

type LLMChain struct {
	*chain
	llm    schema.LLM
	prompt *prompt.Template
	opts   LLMChainOptions
}

func NewLLMChain(llm schema.LLM, prompt *prompt.Template, optFns ...func(o *LLMChainOptions)) (*LLMChain, error) {
	opts := LLMChainOptions{
		OutputKey: "text",
		callbackOptions: callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	llmChain := &LLMChain{
		prompt: prompt,
		llm:    llm,
		opts:   opts,
	}

	llmChain.chain = newChain(llmChain.call, prompt.InputVariables(), []string{opts.OutputKey})

	return llmChain, nil
}

func (c *LLMChain) Type() string {
	return "llm_chain"
}

func (c *LLMChain) Predict(ctx context.Context, values schema.ChainValues) (string, error) {
	output, err := c.Call(ctx, values)
	if err != nil {
		return "", err
	}

	return output[c.opts.OutputKey].(string), err
}

func (c *LLMChain) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	cm := callback.NewManager(c.opts.Callbacks, c.opts.Verbose)

	if err := cm.OnChainStart("LLMChain", &values); err != nil {
		return nil, err
	}

	promptValue, err := c.prompt.FormatPrompt(values)
	if err != nil {
		return nil, err
	}

	res, err := c.llm.GeneratePrompt(ctx, []schema.PromptValue{promptValue}, func(o *schema.GenerateOptions) {
		o.Callbacks = c.opts.Callbacks
	})
	if err != nil {
		return nil, err
	}

	output, err := c.getFinalOutput(res.Generations[0])
	if err != nil {
		return nil, err
	}

	if err := cm.OnChainEnd(&schema.ChainValues{"outputs": output}); err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: output,
	}, nil
}

func (c *LLMChain) Prompt() *prompt.Template {
	return c.prompt
}

func (c *LLMChain) getFinalOutput(generations []*schema.Generation) (any, error) { // nolint unparam
	completion := generations[0].Text
	// TODO Outputparser
	return completion, nil
}
