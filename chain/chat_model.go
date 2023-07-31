package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure ChatModel satisfies the Chain interface.
var _ schema.Chain = (*ChatModel)(nil)

// ChatModelOptions contains options for the ChatModel chain.
type ChatModelOptions struct {
	// CallbackOptions contains options for the chain callbacks.
	*schema.CallbackOptions

	// OutputKey is the key to access the output value containing the ChatModel response summary.
	OutputKey string
}

// ChatModel represents a chain that interacts with a ChatModel and a prompt template.
type ChatModel struct {
	chatModel schema.ChatModel
	prompt    prompt.ChatTemplate
	functions []schema.FunctionDefinition
	opts      ChatModelOptions
}

// NewChatModel creates a new ChatModel chain with the given ChatModel and prompt template.
func NewChatModel(chatModel schema.ChatModel, prompt prompt.ChatTemplate, optFns ...func(o *ChatModelOptions)) (*ChatModel, error) {
	return NewChatModelWithFunctions(chatModel, prompt, nil, optFns...)
}

// NewChatModelWithFunctions creates a new ChatModel chain with the given ChatModel, prompt template, and function definitions.
func NewChatModelWithFunctions(chatModel schema.ChatModel, prompt prompt.ChatTemplate, functions []schema.FunctionDefinition, optFns ...func(o *ChatModelOptions)) (*ChatModel, error) {
	opts := ChatModelOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		OutputKey: "message",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &ChatModel{
		chatModel: chatModel,
		prompt:    prompt,
		functions: functions,
		opts:      opts,
	}, nil
}

// Call executes the ChatModel chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *ChatModel) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	pv, err := c.prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	result, err := model.GeneratePrompt(ctx, c.chatModel, pv, func(o *model.Options) {
		o.Functions = c.functions
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: result.Generations[0].Message,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *ChatModel) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *ChatModel) Type() string {
	return "ChatModel"
}

// Verbose returns the verbosity setting of the chain.
func (c *ChatModel) Verbose() bool {
	return c.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *ChatModel) Callbacks() []schema.Callback {
	return c.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (c *ChatModel) InputKeys() []string {
	return c.prompt.InputVariables()
}

// OutputKeys returns the output keys the chain will return.
func (c *ChatModel) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
