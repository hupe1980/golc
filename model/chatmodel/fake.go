package chatmodel

import (
	"context"

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

func (cm *Fake) Generate(ctx context.Context, messages schema.ChatMessages) (*schema.LLMResult, error) {
	return &schema.LLMResult{
		Generations: [][]*schema.Generation{{newChatGeneraton(cm.response)}},
		LLMOutput:   map[string]any{},
	}, nil
}

func (cm *Fake) Type() string {
	return "Fake"
}

func (cm *Fake) Verbose() bool {
	return false
}

func (cm *Fake) Callbacks() []schema.Callback {
	return []schema.Callback{}
}
