package chain

import (
	"context"
	"errors"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAIModeration satisfies the Chain interface.
var _ schema.Chain = (*OpenAIModeration)(nil)

type OpenAIClient interface {
	Moderations(ctx context.Context, request openai.ModerationRequest) (response openai.ModerationResponse, err error)
}

type OpenAIModerateFunc func(id, model string, result openai.Result) (schema.ChainValues, error)

type OpenAIModerationOptions struct {
	*schema.CallbackOptions
	ModelName          string
	InputKey           string
	OutputKey          string
	OpenAIModerateFunc OpenAIModerateFunc
}

type OpenAIModeration struct {
	client OpenAIClient
	opts   OpenAIModerationOptions
}

func NewOpenAIModeration(apiKey string, optFns ...func(o *OpenAIModerationOptions)) (*OpenAIModeration, error) {
	client := openai.NewClient(apiKey)
	return NewOpenAIModerationFromClient(client, optFns...)
}

func NewOpenAIModerationFromClient(client OpenAIClient, optFns ...func(o *OpenAIModerationOptions)) (*OpenAIModeration, error) {
	opts := OpenAIModerationOptions{
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

	if opts.OpenAIModerateFunc == nil {
		opts.OpenAIModerateFunc = func(id, model string, result openai.Result) (schema.ChainValues, error) {
			if result.Flagged {
				return nil, errors.New("content policy violation")
			}

			return schema.ChainValues{
				opts.OutputKey: result,
			}, nil
		}
	}

	return &OpenAIModeration{
		client: client,
		opts:   opts,
	}, nil
}

// Call executes the openai moderation chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *OpenAIModeration) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	input, ok := inputs[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	text, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	res, err := c.client.Moderations(ctx, openai.ModerationRequest{
		Model: c.opts.ModelName,
		Input: text,
	})
	if err != nil {
		return nil, err
	}

	return c.opts.OpenAIModerateFunc(res.ID, res.Model, res.Results[0])
}

// Memory returns the memory associated with the chain.
func (c *OpenAIModeration) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *OpenAIModeration) Type() string {
	return "OpenAIModeration"
}

// Verbose returns the verbosity setting of the chain.
func (c *OpenAIModeration) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *OpenAIModeration) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *OpenAIModeration) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *OpenAIModeration) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
