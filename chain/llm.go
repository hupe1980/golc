package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure LLM satisfies the Chain interface.
var _ schema.Chain = (*LLM)(nil)

// LLMOptions contains options for the LLM chain.
type LLMOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// Memory is the schema.Memory to be associated with the chain.
	Memory schema.Memory

	// OutputKey is the key to access the output value containing the LLM text generation.
	OutputKey string

	// OutputParser is the schema.OutputParser[any] instance used to parse the LLM text generation result.
	OutputParser schema.OutputParser[any]

	// ReturnFinalOnly determines whether to return only the final parsed result or include extra generation information.
	// When set to true (default), the field will return only the final parsed result.
	// If set to false, the field will include additional information about the generation along with the final parsed result.
	ReturnFinalOnly bool
}

// LLM is a chain implementation that uses the Language Model (LLM) to generate text based on a given prompt.
type LLM struct {
	llm    schema.Model
	prompt *prompt.Template
	opts   LLMOptions
}

// NewLLM creates a new instance of the LLM chain.
func NewLLM(llm schema.Model, prompt *prompt.Template, optFns ...func(o *LLMOptions)) (*LLM, error) {
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

// Call executes the llm chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *LLM) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	promptValue, err := c.prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: fmt.Sprintf("\nPrompt after formatting:\n%s", promptValue.String()),
	}); cbErr != nil {
		return nil, cbErr
	}

	res, err := model.GeneratePrompt(ctx, c.llm, promptValue, func(o *model.Options) {
		o.Stop = opts.Stop
		o.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		o.ParentRunID = opts.CallbackManger.RunID()
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

// GetNumTokens returns the number of tokens in the given text for the associated Language Model (LLM).
func (c *LLM) GetNumTokens(text string) (uint, error) {
	return c.llm.GetNumTokens(text)
}

// Prompt returns the prompt.Template associated with the chain.
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

func (c *LLM) createOutputs(modelResult *schema.ModelResult) ([]map[string]any, error) {
	result := make([]map[string]any, len(modelResult.Generations))

	for i, generation := range modelResult.Generations {
		parsed, err := c.opts.OutputParser.ParseResult(generation)
		if err != nil {
			return nil, err
		}

		result[i] = map[string]any{
			c.opts.OutputKey: parsed,
		}

		if !c.opts.ReturnFinalOnly {
			result[i]["fullGeneration"] = generation
		}
	}

	return result, nil
}
