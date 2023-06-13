package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
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
	InputKey  string
	OutputKey string
}

type LLMBashChain struct {
	chain       *LLMChain
	bashProcess *integration.BashProcess
	opts        LLMBashChainOptions
}

func NewLLMBashChain(chain *LLMChain) (*LLMBashChain, error) {
	opts := LLMBashChainOptions{
		InputKey:  "question",
		OutputKey: "answer",
	}

	bp, err := integration.NewBashProcess()
	if err != nil {
		return nil, err
	}

	return &LLMBashChain{
		chain:       chain,
		bashProcess: bp,
		opts:        opts,
	}, nil
}

func NewLLMBashChainFromLLM(llm golc.LLM) (*LLMBashChain, error) {
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

func (lc *LLMBashChain) Call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	input, ok := values[lc.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, lc.opts.InputKey)
	}

	question, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	t, err := Run(ctx, lc.chain, question)
	if err != nil {
		return nil, err
	}

	outputParser, ok := lc.chain.Prompt().OutputParser()
	if !ok {
		return nil, ErrNoOutputParser
	}

	commands, err := outputParser.Parse(t)
	if err != nil {
		return nil, err
	}

	output, err := lc.bashProcess.Run(commands.([]string))
	if err != nil {
		return nil, err
	}

	return golc.ChainValues{
		lc.opts.OutputKey: output,
	}, nil
}

// InputKeys returns the expected input keys.
func (lc *LLMBashChain) InputKeys() []string {
	return []string{lc.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (lc *LLMBashChain) OutputKeys() []string {
	return []string{lc.opts.OutputKey}
}
