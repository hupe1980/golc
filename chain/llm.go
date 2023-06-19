package chain

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type LLMChainOptions struct {
	*callbackOptions
	Memory       schema.Memory
	OutputKey    string
	OutputParser schema.OutputParser[any]
}

type LLMChain struct {
	*baseChain
	llm    schema.LLM
	prompt *prompt.Template
	opts   LLMChainOptions
}

func NewLLMChain(llm schema.LLM, prompt *prompt.Template, optFns ...func(o *LLMChainOptions)) (*LLMChain, error) {
	opts := LLMChainOptions{
		OutputKey: "text",
		callbackOptions: &callbackOptions{
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

	llmChain.baseChain = &baseChain{
		chainName:       "LLMChain",
		callFunc:        llmChain.call,
		inputKeys:       prompt.InputVariables(),
		outputKeys:      []string{opts.OutputKey},
		memory:          opts.Memory,
		callbackOptions: opts.callbackOptions,
	}

	return llmChain, nil
}

func (c *LLMChain) Predict(ctx context.Context, inputs schema.ChainValues) (string, error) {
	output, err := c.Call(ctx, inputs)
	if err != nil {
		return "", err
	}

	return output[c.opts.OutputKey].(string), err
}

func (c *LLMChain) call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	promptValue, err := c.prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	res, err := c.llm.GeneratePrompt(ctx, []schema.PromptValue{promptValue}, func(o *schema.GenerateOptions) {
		o.Callbacks = c.opts.Callbacks
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: c.getFinalOutput(res.Generations),
	}, nil
}

func (c *LLMChain) Prompt() *prompt.Template {
	return c.prompt
}

func (c *LLMChain) Memory() schema.Memory {
	return c.opts.Memory
}

func (c *LLMChain) Type() string {
	return "LLM"
}

func (c *LLMChain) Verbose() bool {
	return c.opts.callbackOptions.Verbose
}

func (c *LLMChain) Callbacks() []schema.Callback {
	return c.opts.callbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *LLMChain) InputKeys() []string {
	return c.prompt.InputVariables()
}

// OutputKeys returns the output keys the chain will return.
func (c *LLMChain) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *LLMChain) getFinalOutput(generations [][]*schema.Generation) string {
	output := []string{}
	for _, generation := range generations {
		// Get the text of the top generated string.
		output = append(output, strings.TrimSpace(generation[0].Text))
	}

	return output[0]
}
