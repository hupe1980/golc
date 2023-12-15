package chatmodel

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Fake satisfies the ChatModel interface.
var _ schema.ChatModel = (*Fake)(nil)

// FakeResultFunc is a function type used for providing custom model results in the Fake model.
type FakeResultFunc func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error)

// FakeOptions contains options for configuring the Fake model.
type FakeOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
	ChatModelType           string `map:"-"`
}

// Fake is a mock implementation of the schema.ChatModel interface for testing purposes.
type Fake struct {
	schema.Tokenizer
	fakeResultFunc FakeResultFunc
	opts           FakeOptions
}

// NewSimpleFake creates a simple instance of the Fake model with a fixed response for all inputs.
func NewSimpleFake(response string, optFns ...func(o *FakeOptions)) *Fake {
	return NewFake(func(ctx context.Context, messages schema.ChatMessages) (*schema.ModelResult, error) {
		return &schema.ModelResult{
			Generations: []schema.Generation{newChatGeneraton(response)},
			LLMOutput:   map[string]any{},
		}, nil
	}, optFns...)
}

// NewFake creates an instance of the Fake model with the provided custom result function.
func NewFake(fakeResultFunc FakeResultFunc, optFns ...func(o *FakeOptions)) *Fake {
	opts := FakeOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ChatModelType: "chatmodel.Fake",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Fake{
		fakeResultFunc: fakeResultFunc,
		opts:           opts,
	}
}

// Generate generates text based on the provided chat messages and options.
func (cm *Fake) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return cm.fakeResultFunc(ctx, messages)
}

// Type returns the type of the model.
func (cm *Fake) Type() string {
	return cm.opts.ChatModelType
}

// Verbose returns the verbosity setting of the model.
func (cm *Fake) Verbose() bool {
	return cm.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Fake) Callbacks() []schema.Callback {
	return []schema.Callback{}
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Fake) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}
