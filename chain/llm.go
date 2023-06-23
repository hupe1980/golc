package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure LLM satisfies the Chain interface.
var _ schema.Chain = (*LLM)(nil)

type LLMOptions struct {
	*schema.CallbackOptions
	Memory       schema.Memory
	OutputKey    string
	OutputParser schema.OutputParser[any]
	// ReturnFinalOnly determines whether to return only the final parsed result or include extra generation information.
	// When set to true (default), the field will return only the final parsed result.
	// If set to false, the field will include additional information about the generation along with the final parsed result.
	ReturnFinalOnly bool
}

type LLM struct {
	llm    schema.LLM
	prompt *prompt.Template
	opts   LLMOptions
}

func NewLLM(llm schema.LLM, prompt *prompt.Template, optFns ...func(o *LLMOptions)) (*LLM, error) {
	opts := LLMOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		OutputKey:       "text",
		ReturnFinalOnly: true,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.OutputParser == nil {
		opts.OutputParser = outputparser.NewNoOpt()
	}

	return &LLM{
		prompt: prompt,
		llm:    llm,
		opts:   opts,
	}, nil
}

// Call executes the ConversationalRetrieval chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *LLM) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	promptValue, err := c.prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	res, err := model.GeneratePrompt(ctx, c.llm, []schema.PromptValue{promptValue}, func(o *model.Options) {
		o.Stop = opts.Stop

		if opts.CallbackManger != nil {
			o.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
			o.ParentRunID = opts.CallbackManger.RunID()
		}
	})
	if err != nil {
		return nil, err
	}

	outputs, err := c.createOutputs(res)
	if err != nil {
		return nil, err
	}

	return outputs[0], nil
}

func (c *LLM) Prompt() *prompt.Template {
	return c.prompt
}

// Memory returns the memory associated with the chain.
func (c *LLM) Memory() schema.Memory {
	return c.opts.Memory
}

// Type returns the type of the chain.
func (c *LLM) Type() string {
	return "LLM"
}

// Verbose returns the verbosity setting of the chain.
func (c *LLM) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *LLM) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *LLM) InputKeys() []string {
	return c.prompt.InputVariables()
}

// OutputKeys returns the output keys the chain will return.
func (c *LLM) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *LLM) createOutputs(llmResult *schema.LLMResult) ([]map[string]any, error) {
	result := make([]map[string]any, len(llmResult.Generations)-1)

	for _, generation := range llmResult.Generations {
		parsed, err := c.opts.OutputParser.ParseResult(generation)
		if err != nil {
			return nil, err
		}

		output := map[string]any{
			c.opts.OutputKey: parsed,
			"fullGeneration": generation,
		}

		result = append(result, output)
	}

	if c.opts.ReturnFinalOnly {
		for i := range result {
			result[i] = map[string]any{
				c.opts.OutputKey: result[i][c.opts.OutputKey],
			}
		}
	}

	return result, nil
}
