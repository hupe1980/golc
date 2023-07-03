package chatmodel

import (
	"context"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Fake satisfies the ChatModel interface.
var _ schema.ChatModel = (*Fake)(nil)

type Fake struct {
	schema.Tokenizer
	response string
}

func NewFake(response string) *Fake {
	return &Fake{
		response: response,
	}
}

func (cm *Fake) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &schema.ModelResult{
		Generations: [][]schema.Generation{{newChatGeneraton(cm.response)}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (cm *Fake) Type() string {
	return "chatmodel.Fake"
}

func (cm *Fake) Verbose() bool {
	return false
}

func (cm *Fake) Callbacks() []schema.Callback {
	return []schema.Callback{}
}

func (cm *Fake) InvocationParams() map[string]any {
	return nil
}
