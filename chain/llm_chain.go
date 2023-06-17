package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/prompt"
)

type LLMChainOptions struct {
	callbackOptions
	OutputKey    string
	OutputParser golc.OutputParser[any]
}

type LLMChain struct {
	llm    golc.LLM
	prompt *prompt.Template
	opts   LLMChainOptions
}

func NewLLMChain(llm golc.LLM, prompt *prompt.Template, optFns ...func(o *LLMChainOptions)) (*LLMChain, error) {
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

	return llmChain, nil
}

func (c *LLMChain) Type() string {
	return "llm_chain"
}

func (c *LLMChain) Predict(ctx context.Context, values golc.ChainValues) (string, error) {
	output, err := c.Call(ctx, values)
	if err != nil {
		return "", err
	}

	return output[c.opts.OutputKey].(string), err
}

func (c *LLMChain) Call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	cm := callback.NewManager(c.opts.Callbacks, c.opts.Verbose)

	if err := cm.OnChainStart("LLMChain", &values); err != nil {
		return nil, err
	}

	promptValue, err := c.prompt.FormatPrompt(values)
	if err != nil {
		return nil, err
	}

	res, err := c.llm.GeneratePrompt(ctx, []golc.PromptValue{promptValue}, func(o *golc.GenerateOptions) {
		o.Callbacks = c.opts.Callbacks
	})
	if err != nil {
		return nil, err
	}

	output, err := c.getFinalOutput(res.Generations[0])
	if err != nil {
		return nil, err
	}

	if err := cm.OnChainEnd(&golc.ChainValues{"outputs": output}); err != nil {
		return nil, err
	}

	return golc.ChainValues{
		c.opts.OutputKey: output,
	}, nil
}

// InputKeys returns the expected input keys.
func (c *LLMChain) InputKeys() []string {
	return append([]string{}, c.prompt.InputVariables()...)
}

// OutputKeys returns the output keys the chain will return.
func (c *LLMChain) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *LLMChain) Prompt() *prompt.Template {
	return c.prompt
}

func (c *LLMChain) getFinalOutput(generations []*golc.Generation) (any, error) { // nolint unparam
	completion := generations[0].Text
	// TODO Outputparser
	return completion, nil
}
