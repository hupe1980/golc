package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const llmBashTemplate = `If someone asks you to perform a task, your job is to come up with a series of bash commands that will perform the task. There is no need to put "#!/bin/bash" in your answer. Make sure to reason step by step, using this format:

Question: "copy the files in the directory named 'target' into a new directory at the same level as target called 'myNewDirectory'"

I need to take the following actions:
- List all files in the directory
- Create a new directory
- Copy the files from the first directory into the second directory
` + "```" + `bash
ls
mkdir myNewDirectory
cp -r target/* myNewDirectory
` + "```" + `

That is the format. Begin!

Question: {{.question}}`

type LLMBashOptions struct {
	*callbackOptions
	InputKey  string
	OutputKey string
}

type LLMBash struct {
	*baseChain
	llmChain    *LLMChain
	bashProcess *integration.BashProcess
	opts        LLMBashOptions
}

func NewLLMBash(llmChain *LLMChain, optFns ...func(o *LLMBashOptions)) (*LLMBash, error) {
	opts := LLMBashOptions{
		InputKey:  "question",
		OutputKey: "answer",
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	bp, err := integration.NewBashProcess()
	if err != nil {
		return nil, err
	}

	bash := &LLMBash{
		llmChain:    llmChain,
		bashProcess: bp,
		opts:        opts,
	}

	bash.baseChain = &baseChain{
		chainName:       "LLMBash",
		callFunc:        bash.call,
		inputKeys:       []string{opts.InputKey},
		outputKeys:      []string{opts.OutputKey},
		callbackOptions: opts.callbackOptions,
	}

	return bash, nil
}

func NewLLMBashFromLLM(llm schema.LLM) (*LLMBash, error) {
	prompt, err := prompt.NewTemplate(llmBashTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewFencedCodeBlock("```bash")
	})
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, prompt)
	if err != nil {
		return nil, err
	}

	return NewLLMBash(llmChain)
}

func (c *LLMBash) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	question, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	t, err := c.llmChain.Run(ctx, question)
	if err != nil {
		return nil, err
	}

	outputParser, ok := c.llmChain.Prompt().OutputParser()
	if !ok {
		return nil, ErrNoOutputParser
	}

	commands, err := outputParser.Parse(t)
	if err != nil {
		return nil, err
	}

	output, err := c.bashProcess.Run(ctx, commands.([]string))
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: output,
	}, nil
}

func (c *LLMBash) Memory() schema.Memory {
	return nil
}

func (c *LLMBash) Type() string {
	return "LLMBash"
}

func (c *LLMBash) Verbose() bool {
	return c.opts.callbackOptions.Verbose
}

func (c *LLMBash) Callbacks() []schema.Callback {
	return c.opts.callbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *LLMBash) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *LLMBash) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
