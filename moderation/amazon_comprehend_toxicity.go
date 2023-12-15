package moderation

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// AmazonComprehendToxicityClient is an interface for the Amazon Comprehend client used for toxicity detection.
type AmazonComprehendToxicityClient interface {
	// DetectToxicContent is an interface function that analyzes a given text for toxic content.
	DetectToxicContent(ctx context.Context, params *comprehend.DetectToxicContentInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectToxicContentOutput, error)
}

// AmazonComprehendToxicityOptions contains options for configuring the AmazonComprehendToxicity instance.
type AmazonComprehendToxicityOptions struct {
	// CallbackOptions embeds CallbackOptions to include the verbosity setting and callbacks.
	*schema.CallbackOptions
	// InputKey is the key to extract the input text from the input ChainValues.
	InputKey string
	// OutputKey is the key to store the output of the moderation in the output ChainValues.
	OutputKey string
	// LanguageCode is the language code used for toxicity detection (default is "en").
	LanguageCode string
	// Labels is a list of specific toxicity labels to check for.
	Labels []string
	// Threshold is the threshold score for considering content as toxic (default is 0.8).
	Threshold float32
}

// AmazonComprehendToxicity is a content moderation chain using Amazon Comprehend for toxicity detection.
type AmazonComprehendToxicity struct {
	client AmazonComprehendToxicityClient
	opts   AmazonComprehendToxicityOptions
}

// NewAmazonComprehendToxicity creates a new instance of AmazonComprehendToxicity with the provided client and options.
func NewAmazonComprehendToxicity(client AmazonComprehendToxicityClient, optFns ...func(o *AmazonComprehendToxicityOptions)) *AmazonComprehendToxicity {
	opts := AmazonComprehendToxicityOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:     "input",
		OutputKey:    "output",
		LanguageCode: "en",
		Threshold:    0.8,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonComprehendToxicity{
		client: client,
		opts:   opts,
	}
}

// Call executes the amazon comprehend moderation chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *AmazonComprehendToxicity) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
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

	// TODO split text

	output, err := c.client.DetectToxicContent(ctx, &comprehend.DetectToxicContentInput{
		TextSegments: []types.TextSegment{{Text: aws.String(text)}},
		LanguageCode: types.LanguageCode(c.opts.LanguageCode),
	})
	if err != nil {
		return nil, err
	}

	if len(c.opts.Labels) == 0 {
		for _, item := range output.ResultList {
			if aws.ToFloat32(item.Toxicity) >= c.opts.Threshold {
				return nil, errors.New("toxic content found")
			}
		}
	} else {
		for _, item := range output.ResultList {
			for _, label := range item.Labels {
				if util.Contains(c.opts.Labels, string(label.Name)) {
					if aws.ToFloat32(label.Score) >= c.opts.Threshold {
						return nil, errors.New("toxic content found")
					}
				}
			}
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: text,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *AmazonComprehendToxicity) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *AmazonComprehendToxicity) Type() string {
	return "AmazonComprehendToxicityModeration"
}

// Verbose returns the verbosity setting of the chain.
func (c *AmazonComprehendToxicity) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *AmazonComprehendToxicity) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *AmazonComprehendToxicity) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *AmazonComprehendToxicity) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
