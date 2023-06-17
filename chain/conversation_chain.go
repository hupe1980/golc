package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type ConversationChainOptions struct {
	*callbackOptions
	Memory       schema.Memory
	OutputParser schema.OutputParser[any]
}

type ConversationChain struct {
	*baseChain
	llm    schema.LLM
	prompt *prompt.Template
	opts   ConversationChainOptions
}

func NewConversationChain(llm schema.LLM, prompt *prompt.Template, optFns ...func(o *ConversationChainOptions)) (*ConversationChain, error) {
	opts := ConversationChainOptions{
		Memory: memory.NewConversationBuffer(),
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	conversationChain := &ConversationChain{
		prompt: prompt,
		llm:    llm,
		opts:   opts,
	}

	conversationChain.baseChain = &baseChain{
		chainName:       "ConversationChain",
		callFunc:        conversationChain.call,
		inputKeys:       []string{"input"},
		outputKeys:      []string{"response"},
		memory:          opts.Memory,
		callbackOptions: opts.callbackOptions,
	}

	return conversationChain, nil
}

func (c *ConversationChain) call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return nil, nil
}
