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

const llmBashChainPromptTemplate = `If someone asks you to perform a task, your job is to come up with a series of bash commands that will perform the task. There is no need to put "#!/bin/bash" in your answer. Make sure to reason step by step, using this format:

Question: "copy the files in the directory named 'target' into a new directory at the same level as target called 'myNewDirectory'"

I need to take the following actions:
- List all files in the directory
- Create a new directory
- Copy the files from the first directory into the second directory
` +
	"```bash" + `
ls
mkdir myNewDirectory
cp -r target/* myNewDirectory
` +
	"```" + `

That is the format. Begin!

Question: {{.question}}`

type LLMBashChainOptions struct {
	*callbackOptions
	InputKey  string
	OutputKey string
}

type LLMBashChain struct {
	*baseChain
	llmChain    *LLMChain
	bashProcess *integration.BashProcess
	opts        LLMBashChainOptions
}

func NewLLMBashChain(llmChain *LLMChain, optFns ...func(o *LLMBashChainOptions)) (*LLMBashChain, error) {
	opts := LLMBashChainOptions{
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

	bash := &LLMBashChain{
		llmChain:    llmChain,
		bashProcess: bp,
		opts:        opts,
	}

	bash.baseChain = &baseChain{
		chainName:       "LLMBashChain",
		callFunc:        bash.call,
		inputKeys:       []string{opts.InputKey},
		outputKeys:      []string{opts.OutputKey},
		callbackOptions: opts.callbackOptions,
	}

	return bash, nil
}

func NewLLMBashChainFromLLM(llm schema.LLM) (*LLMBashChain, error) {
	prompt, err := prompt.NewTemplate(llmBashChainPromptTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewBashOutputParser()
	})
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, prompt)
	if err != nil {
		return nil, err
	}

	return NewLLMBashChain(llmChain)
}

func (lc *LLMBashChain) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[lc.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, lc.opts.InputKey)
	}

	question, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	t, err := lc.llmChain.Run(ctx, question)
	if err != nil {
		return nil, err
	}

	outputParser, ok := lc.llmChain.Prompt().OutputParser()
	if !ok {
		return nil, ErrNoOutputParser
	}

	commands, err := outputParser.Parse(t)
	if err != nil {
		return nil, err
	}

	output, err := lc.bashProcess.Run(ctx, commands.([]string))
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		lc.opts.OutputKey: output,
	}, nil
}
