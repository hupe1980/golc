package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const llmMathTemplate = `Translate a math problem into an expression that can be executed using golangs expr library.
Place the expression between a fenced block of code that starts with ` + "```" + `text and ends with ` + "```" + `.
Use the output of running this code to answer the question.

Question: $(Question with math problem)
` + "```" + `text
$(single line mathematical expression that solves the problem)
` + "```" + `

Begin.

Question: What is 37593 * 67?
` + "```" + `text
37593 * 67
` + "```" + `

Question: 37593^(1/5)
` + "```" + `text
37593**(1/5)
` + "```" + `

Question: {{.question}}
`

// Compile time check to ensure LLMMath satisfies the Chain interface.
var _ schema.Chain = (*LLMMath)(nil)

type LLMMathOptions struct {
	*schema.CallbackOptions
	InputKey  string
	OutputKey string
}

type LLMMath struct {
	llmChain *LLM
	opts     LLMMathOptions
}

func NewLLMMath(llm schema.LLM, optFns ...func(o *LLMMathOptions)) (*LLMMath, error) {
	opts := LLMMathOptions{
		InputKey:  "question",
		OutputKey: "answer",
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	prompt := prompt.NewTemplate(llmMathTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewFencedCodeBlock("```text")
	})

	llmChain, err := NewLLM(llm, prompt)
	if err != nil {
		return nil, err
	}

	return &LLMMath{
		llmChain: llmChain,
		opts:     opts,
	}, nil
}

// Call executes the ConversationalRetrieval chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *LLMMath) Call(ctx context.Context, values schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
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

	parsed, err := outputParser.Parse(strings.TrimSpace(t))
	if err != nil {
		return nil, err
	}

	if len(parsed.([]string)) != 1 {
		return nil, fmt.Errorf("unknown format from LLM: %s", t)
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nExpression:\n%s", parsed),
	}); cbErr != nil {
		return nil, cbErr
	}

	output, err := c.evaluateExpression(parsed.([]string)[0])
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

func (c *LLMMath) evaluateExpression(expression string) (string, error) {
	output, err := expr.Eval(expression, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%f", output), nil
}

// Memory returns the memory associated with the chain.
func (c *LLMMath) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *LLMMath) Type() string {
	return "LLMMath"
}

// Verbose returns the verbosity setting of the chain.
func (c *LLMMath) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *LLMMath) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *LLMMath) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *LLMMath) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
