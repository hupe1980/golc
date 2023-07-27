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

// BashRunner is the interface used to execute Bash commands.
type BashRunner interface {
	Run(ctx context.Context, commands []string) (string, error)
}

// VerifyCommands is a function signature used to verify the validity of the
// generated Bash commands before execution.
type VerifyCommands func(commands []string) bool

// BashOptions contains options for the Bash chain.
type BashOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// InputKey is the key to access the input value containing the user question.
	InputKey string

	// OutputKey is the key to access the output value containing the Bash commands output.
	OutputKey string

	// BashRunner is the BashRunner instance used to execute the generated Bash commands.
	BashRunner BashRunner

	// VerifyCommands is a function used to verify the validity of the generated Bash commands before execution.
	// It should return true if the commands are valid, false otherwise.
	VerifyCommands VerifyCommands
}

// Bash is a chain implementation that prompts the user to provide a series of Bash commands
// to perform a specific task based on a given question. It then verifies and executes the provided commands.
type Bash struct {
	llmChain *LLM
	opts     BashOptions
}

// NewBash creates a new instance of the Bash chain.
func NewBash(llm schema.LLM, optFns ...func(o *BashOptions)) (*Bash, error) {
	opts := BashOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:       "question",
		OutputKey:      "answer",
		BashRunner:     integration.NewBashProcess(),
		VerifyCommands: func(commands []string) bool { return true },
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

	return &Bash{
		llmChain: llmChain,
		opts:     opts,
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

	commandStrs, ok := commands.([]string)
	if !ok {
		return nil, fmt.Errorf("cannot convert commands to string: %s", commands)
	}

	if ok := c.opts.VerifyCommands(commandStrs); !ok {
		return nil, fmt.Errorf("invalid commands: %s", commandStrs)
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nCode:\n%s", commandStrs),
	}); cbErr != nil {
		return nil, cbErr
	}

	output, err := c.opts.BashRunner.Run(ctx, commandStrs)
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
