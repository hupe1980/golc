package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultBashTemplate = `If someone asks you to perform a task, your job is to come up with a series of bash commands that will perform the task. There is no need to put "#!/bin/bash" in your answer. Make sure to reason step by step, using this format:

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

// Compile time check to ensure Bash satisfies the Chain interface.
var _ schema.Chain = (*Bash)(nil)

type BashOptions struct {
	*schema.CallbackOptions
	InputKey  string
	OutputKey string
}

type Bash struct {
	llmChain    *LLM
	bashProcess *integration.BashProcess
	opts        BashOptions
}

func NewBash(llm schema.LLM, optFns ...func(o *BashOptions)) (*Bash, error) {
	opts := BashOptions{
		InputKey:  "question",
		OutputKey: "answer",
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	prompt := prompt.NewTemplate(defaultBashTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewFencedCodeBlock("```bash")
	})

	llmChain, err := NewLLM(llm, prompt)
	if err != nil {
		return nil, err
	}

	bp, err := integration.NewBashProcess()
	if err != nil {
		return nil, err
	}

	return &Bash{
		llmChain:    llmChain,
		bashProcess: bp,
		opts:        opts,
	}, nil
}

// Call executes the bash chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *Bash) Call(ctx context.Context, values schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	input, ok := values[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	question, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: question,
	}); cbErr != nil {
		return nil, cbErr
	}

	t, err := golc.SimpleCall(ctx, c.llmChain, question, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
	})
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

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nCode:\n%s", commands),
	}); cbErr != nil {
		return nil, cbErr
	}

	output, err := c.bashProcess.Run(ctx, commands.([]string))
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nAnswer:\n%s", output),
	}); cbErr != nil {
		return nil, cbErr
	}

	return schema.ChainValues{
		c.opts.OutputKey: output,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *Bash) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *Bash) Type() string {
	return "Bash"
}

// Verbose returns the verbosity setting of the chain.
func (c *Bash) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated chain.
func (c *Bash) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Bash) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *Bash) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
