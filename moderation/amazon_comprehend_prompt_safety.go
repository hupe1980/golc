package moderation

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

// AmazonComprehendPromptSafetyClient is an interface for the Amazon Comprehend Prompt Safety client.
type AmazonComprehendPromptSafetyClient interface {
	// Options returns the options associated with the Amazon Comprehend Prompt Safety client.
	Options() comprehend.Options
	// ClassifyDocument analyzes the provided text and classifies it as safe or unsafe based on predefined categories.
	ClassifyDocument(ctx context.Context, params *comprehend.ClassifyDocumentInput, optFns ...func(*comprehend.Options)) (*comprehend.ClassifyDocumentOutput, error)
}

// AmazonComprehendPromptSafetyOptions contains options for the Amazon Comprehend Prompt Safety moderation.
type AmazonComprehendPromptSafetyOptions struct {
	// CallbackOptions embeds CallbackOptions to include the verbosity setting and callbacks.
	*schema.CallbackOptions
	// InputKey is the key to extract the input text from the input ChainValues.
	InputKey string
	// OutputKey is the key to store the output of the moderation in the output ChainValues.
	OutputKey string
	// Threshold is the confidence threshold for determining if an input is considered unsafe.
	Threshold float32
	// Endpoint is the URL endpoint for the external service that performs unsafe content detection.
	Endpoint string
}

// AmazonComprehendPromptSafety is a struct representing the Amazon Comprehend Prompt Safety moderation functionality.
type AmazonComprehendPromptSafety struct {
	client AmazonComprehendPromptSafetyClient
	opts   AmazonComprehendPromptSafetyOptions
}

// NewAmazonComprehendPromptSafety creates a new instance of AmazonComprehendPromptSafety with the provided client and options.
func NewAmazonComprehendPromptSafety(client AmazonComprehendPromptSafetyClient, optFns ...func(o *AmazonComprehendPromptSafetyOptions)) *AmazonComprehendPromptSafety {
	opts := AmazonComprehendPromptSafetyOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:  "input",
		OutputKey: "output",
		Threshold: 0.8,
		Endpoint:  "document-classifier-endpoint/prompt-safety",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonComprehendPromptSafety{
		client: client,
		opts:   opts,
	}
}

// Call executes the amazon comprehend moderation chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *AmazonComprehendPromptSafety) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
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

	output, err := c.client.ClassifyDocument(ctx, &comprehend.ClassifyDocumentInput{
		Text:        aws.String(text),
		EndpointArn: aws.String(fmt.Sprintf("arn:aws:comprehend:%s:aws:%s", c.client.Options().Region, c.opts.Endpoint)),
	})
	if err != nil {
		return nil, err
	}

	for _, classes := range output.Classes {
		if aws.ToString(classes.Name) == "UNSAFE_PROMPT" && aws.ToFloat32(classes.Score) > c.opts.Threshold {
			return nil, errors.New("unsafe prompt detected")
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: text,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *AmazonComprehendPromptSafety) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *AmazonComprehendPromptSafety) Type() string {
	return "AmazonComprehendPIIModeration"
}

// Verbose returns the verbosity setting of the chain.
func (c *AmazonComprehendPromptSafety) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *AmazonComprehendPromptSafety) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *AmazonComprehendPromptSafety) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *AmazonComprehendPromptSafety) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
