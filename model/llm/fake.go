package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Fake satisfies the LLM interface.
var _ schema.LLM = (*Fake)(nil)

// FakeResponseFunc is a function type for generating fake responses based on a prompt.
type FakeResponseFunc func(prompt string) string

// FakeOptions contains options for configuring the Fake LLM model.
type FakeOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
}

// Fake is a fake LLM model that generates text based on a provided response function.
type Fake struct {
	schema.Tokenizer
	responseFunc FakeResponseFunc
	opts         FakeOptions
}

// NewFake creates a new instance of the Fake LLM model with the provided response function and options.
func NewFake(responseFunc FakeResponseFunc, optFns ...func(o *FakeOptions)) *Fake {
	opts := FakeOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Fake{
		responseFunc: responseFunc,
		opts:         opts,
	}
}

// Generate generates text based on the provided prompt and options.
func (l *Fake) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	text := l.responseFunc(prompt)

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: text}},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (l *Fake) Type() string {
	return "llm.Fake"
}

// Verbose returns the verbosity setting of the model.
func (l *Fake) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *Fake) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *Fake) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
