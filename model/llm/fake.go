package llm

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Fake satisfies the LLM interface.
var _ schema.LLM = (*Fake)(nil)

type FakeOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
}

type Fake struct {
	schema.Tokenizer
	response string
}

func NewFake(response string) *Fake {
	return &Fake{
		response: response,
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

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: l.response}},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (l *Fake) Type() string {
	return "llm.Fake"
}

// Verbose returns the verbosity setting of the model.
func (l *Fake) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the model.
func (l *Fake) Callbacks() []schema.Callback {
	return []schema.Callback{}
}

// InvocationParams returns the parameters used in the model invocation.
func (l *Fake) InvocationParams() map[string]any {
	return nil
}
