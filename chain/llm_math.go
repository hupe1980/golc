package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/hupe1980/golc"
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

type LLMMathOptions struct {
	*schema.CallbackOptions
	InputKey  string
	OutputKey string
}

type LLMMath struct {
	llmChain *LLMChain
	opts     LLMMathOptions
}

func NewLLMMath(llmChain *LLMChain, optFns ...func(o *LLMMathOptions)) (*LLMMath, error) {
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

	return &LLMMath{
		llmChain: llmChain,
		opts:     opts,
	}, nil
}

func NewLLMMathFromLLM(llm schema.LLM) (*LLMMath, error) {
	prompt, err := prompt.NewTemplate(llmMathTemplate, func(o *prompt.TemplateOptions) {
		o.OutputParser = outputparser.NewFencedCodeBlock("```text")
	})
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, prompt)
	if err != nil {
		return nil, err
	}

	return NewLLMMath(llmChain)
}

func (c *LLMMath) Call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	question, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	t, err := golc.SimpleCall(ctx, c.llmChain, question)
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

	output, err := c.evaluateExpression(parsed.([]string)[0])
	if err != nil {
		return nil, err
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

func (c *LLMMath) Memory() schema.Memory {
	return nil
}

func (c *LLMMath) Type() string {
	return "LLMMath"
}

func (c *LLMMath) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

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
