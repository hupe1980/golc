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

// defaultMathTemplate defines the default template for translating math problems
// into expressions that can be executed using Golang's expr library.
const defaultMathTemplate = `Translate a math problem into an expression that can be executed using golangs expr library.
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

// Compile time check to ensure Math satisfies the Chain interface.
var _ schema.Chain = (*Math)(nil)

// MathOptions contains options for the Math chain.
type MathOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// InputKey is the key to access the input value containing the user question.
	InputKey string

	// OutputKey is the key to access the output value containing the math expression result.
	OutputKey string
}

// Math is a chain implementation that prompts the user to provide a math problem
// in the form of a single-line mathematical expression that can be executed using
// a golang expr library. It then translates the expression and evaluates it to provide the result.
type Math struct {
	llmChain *LLM
	opts     MathOptions
}

// NewMath creates a new instance of the Math chain.
func NewMath(model schema.Model, optFns ...func(o *MathOptions)) (*Math, error) {
	opts := MathOptions{
		InputKey:  "question",
		OutputKey: "answer",
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	prompt := prompt.NewTemplate(defaultMathTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewFencedCodeBlock("```text")
	})

	llmChain, err := NewLLM(model, prompt)
	if err != nil {
		return nil, err
	}

	return &Math{
		llmChain: llmChain,
		opts:     opts,
	}, nil
}

// Call executes the math chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *Math) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	question, err := inputs.GetString(c.opts.InputKey)
	if err != nil {
		return nil, err
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

// evaluateExpression evaluates the mathematical expression using a golang expr library.
func (c *Math) evaluateExpression(expression string) (string, error) {
	output, err := expr.Eval(expression, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", output), nil
}

// Memory returns the memory associated with the chain.
func (c *Math) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *Math) Type() string {
	return "Math"
}

// Verbose returns the verbosity setting of the chain.
func (c *Math) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *Math) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Math) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *Math) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
