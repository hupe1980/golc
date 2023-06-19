package chain

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const conversationTemplate = `The following is a friendly conversation between a human and an AI. The AI is talkative and provides lots of specific details from its context. If the AI does not know the answer to a question, it truthfully says it does not know.

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

	return &Conversation{
		llm:  llm,
		opts: opts,
	}, nil
}

func (c *Conversation) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
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

func (c *Conversation) Memory() schema.Memory {
	return c.opts.Memory
}

func (c *Conversation) Type() string {
	return "Conversation"
}

func (c *Conversation) Verbose() bool {
	return c.opts.callbackOptions.Verbose
}

func (c *Conversation) Callbacks() []schema.Callback {
	return c.opts.callbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Conversation) InputKeys() []string {
	return []string{"input"}
}

// OutputKeys returns the output keys the chain will return.
func (c *Conversation) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *Conversation) getFinalOutput(generations [][]*schema.Generation) string {
	output := []string{}
	for _, generation := range generations {
		// Get the text of the top generated string.
		output = append(output, strings.TrimSpace(generation[0].Text))
	}

	return output[0]
}
