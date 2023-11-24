package moderation

import (
	"context"
	"errors"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the Chain interface.
var _ schema.Chain = (*OpenAI)(nil)

// OpenAIClient is an interface representing an OpenAI client that can make moderation requests.
type OpenAIClient interface {
	// Moderations sends a moderation request to the OpenAI API and receives the response.
	// It takes the context and a ModerationRequest as input and returns a ModerationResponse or an error.
	Moderations(ctx context.Context, request openai.ModerationRequest) (response openai.ModerationResponse, err error)
}

// OpenAIOptions contains options for configuring the OpenAI chain.
type OpenAIOptions struct {
	// CallbackOptions embeds CallbackOptions to include the verbosity setting and callbacks.
	*schema.CallbackOptions
	// ModelName is the name of the OpenAI model to use for moderation.
	ModelName string
	// InputKey is the key to extract the input text from the input ChainValues.
	InputKey string
	// OutputKey is the key to store the output of the moderation in the output ChainValues.
	OutputKey string
}

// OpenAI represents a chain that performs moderation using the OpenAI API.
type OpenAI struct {
	client OpenAIClient
	opts   OpenAIOptions
}

// NewOpenAI creates a new instance of the OpenAI chain using the provided API key and options.
func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) *OpenAI {
	client := openai.NewClient(apiKey)
	return NewOpenAIFromClient(client, optFns...)
}

// NewOpenAIFromClient creates a new instance of the OpenAI chain with the given OpenAI client and options.
func NewOpenAIFromClient(client OpenAIClient, optFns ...func(o *OpenAIOptions)) *OpenAI {
	opts := OpenAIOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName: "text-moderation-latest",
		InputKey:  "input",
		OutputKey: "output",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &OpenAI{
		client: client,
		opts:   opts,
	}
}

// Call executes the openai moderation chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *OpenAI) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	text, err := inputs.GetString(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
		Text: text,
	}); cbErr != nil {
		return nil, cbErr
	}

	res, err := c.client.Moderations(ctx, openai.ModerationRequest{
		Model: c.opts.ModelName,
		Input: text,
	})
	if err != nil {
		return nil, err
	}

	if res.Results[0].Flagged {
		return nil, errors.New("content policy violation")
	}

	return schema.ChainValues{
		c.opts.OutputKey: text,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *OpenAI) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *OpenAI) Type() string {
	return "OpenAIModeration"
}

// Verbose returns the verbosity setting of the chain.
func (c *OpenAI) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *OpenAI) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *OpenAI) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *OpenAI) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
