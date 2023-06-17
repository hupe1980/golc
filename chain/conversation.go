package chain

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const conversationTemplate = `he following is a friendly conversation between a human and an AI. The AI is talkative and provides lots of specific details from its context. If the AI does not know the answer to a question, it truthfully says it does not know.

Current conversation:
{{.history}}
Human: {{.input}}
AI:`

type ConversationOptions struct {
	*callbackOptions
	Prompt       *prompt.Template
	Memory       schema.Memory
	OutputKey    string
	OutputParser schema.OutputParser[any]
}

type Conversation struct {
	*baseChain
	llm  schema.LLM
	opts ConversationOptions
}

func NewConversation(llm schema.LLM, optFns ...func(o *ConversationOptions)) (*Conversation, error) {
	opts := ConversationOptions{
		OutputKey: "response",
		Memory:    memory.NewConversationBuffer(),
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Prompt == nil {
		var pErr error

		opts.Prompt, pErr = prompt.NewTemplate(conversationTemplate)
		if pErr != nil {
			return nil, pErr
		}
	}

	conversation := &Conversation{
		llm:  llm,
		opts: opts,
	}

	conversation.baseChain = &baseChain{
		chainName:       "Conversation",
		callFunc:        conversation.call,
		inputKeys:       []string{"input"},
		outputKeys:      []string{opts.OutputKey},
		memory:          opts.Memory,
		callbackOptions: opts.callbackOptions,
	}

	return conversation, nil
}

func (c *Conversation) Predict(ctx context.Context, inputs schema.ChainValues) (string, error) {
	output, err := c.Call(ctx, inputs)
	if err != nil {
		return "", err
	}

	return output[c.opts.OutputKey].(string), err
}

func (c *Conversation) call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	promptValue, err := c.opts.Prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	res, err := c.llm.GeneratePrompt(ctx, []schema.PromptValue{promptValue}, func(o *schema.GenerateOptions) {
		o.Callbacks = c.opts.Callbacks
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: c.getFinalOutput(res.Generations),
	}, nil
}

func (c *Conversation) Prompt() *prompt.Template {
	return c.opts.Prompt
}

func (c *Conversation) getFinalOutput(generations [][]*schema.Generation) string {
	output := []string{}
	for _, generation := range generations {
		// Get the text of the top generated string.
		output = append(output, strings.TrimSpace(generation[0].Text))
	}

	return output[0]
}
