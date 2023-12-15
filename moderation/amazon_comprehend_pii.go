package moderation

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// AmazonComprehendPIIClient is an interface for the Amazon Comprehend client.
type AmazonComprehendPIIClient interface {
	// ContainsPiiEntities is an interface method that checks if the input text contains Personally Identifiable Information (PII) entities.
	ContainsPiiEntities(ctx context.Context, params *comprehend.ContainsPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.ContainsPiiEntitiesOutput, error)
	// DetectPiiEntities is an interface method that checks if the input text contains Personally Identifiable Information (PII) entities and returns information about them.
	DetectPiiEntities(ctx context.Context, params *comprehend.DetectPiiEntitiesInput, optFns ...func(*comprehend.Options)) (*comprehend.DetectPiiEntitiesOutput, error)
}

// ReactFunc is a function type that defines how to react to PII entities found in the text.
type RedactFunc func(ctx context.Context, text string, maskMarker rune, entityType string, offsetBegin, offsetEnd int32) string

// AmazonComprehendPIIOptions contains options for the Amazon Comprehend PII moderation.
type AmazonComprehendPIIOptions struct {
	// CallbackOptions embeds CallbackOptions to include the verbosity setting and callbacks.
	*schema.CallbackOptions
	// InputKey is the key to extract the input text from the input ChainValues.
	InputKey string
	// OutputKey is the key to store the output of the moderation in the output ChainValues.
	OutputKey string
	// LanguageCode is the language code to specify the language of the input text.
	LanguageCode string
	// Labels is a list of labels to check for in the PII analysis.
	Labels []string
	// Threshold is the threshold for determining if PII content is found.
	Threshold float32
	// Redact enables redaction of detected PII entities.
	Redact bool
	// MaskMarker is the redaction mask character in case redaction
	MaskMarker rune
	// RedactFunc defines how to redact PII entities found in the text.
	RedactFunc RedactFunc
}

// AmazonComprehendPII is a struct representing the Amazon Comprehend PII moderation functionality.
type AmazonComprehendPII struct {
	client AmazonComprehendPIIClient
	opts   AmazonComprehendPIIOptions
}

// NewAmazonComprehendPII creates a new instance of AmazonComprehendPII with the provided client and options.
func NewAmazonComprehendPII(client AmazonComprehendPIIClient, optFns ...func(o *AmazonComprehendPIIOptions)) *AmazonComprehendPII {
	opts := AmazonComprehendPIIOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:     "input",
		OutputKey:    "output",
		LanguageCode: "en",
		Threshold:    0.8,
		Redact:       false,
		MaskMarker:   '*',
		RedactFunc: func(ctx context.Context, text string, maskMarker rune, entityType string, offsetBegin, offsetEnd int32) string {
			maskLength := offsetEnd - offsetBegin
			maskedPart := strings.Repeat(string(maskMarker), int(maskLength))

			return text[:offsetBegin] + maskedPart + text[offsetEnd:]
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AmazonComprehendPII{
		client: client,
		opts:   opts,
	}
}

// Call executes the amazon comprehend moderation chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *AmazonComprehendPII) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
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

	if !c.opts.Redact {
		return c.containsPII(ctx, text)
	}

	return c.detectPII(ctx, text)
}

// Memory returns the memory associated with the chain.
func (c *AmazonComprehendPII) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *AmazonComprehendPII) Type() string {
	return "AmazonComprehendPIIModeration"
}

// Verbose returns the verbosity setting of the chain.
func (c *AmazonComprehendPII) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *AmazonComprehendPII) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *AmazonComprehendPII) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *AmazonComprehendPII) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *AmazonComprehendPII) containsPII(ctx context.Context, text string) (schema.ChainValues, error) {
	output, err := c.client.ContainsPiiEntities(ctx, &comprehend.ContainsPiiEntitiesInput{
		Text:         aws.String(text),
		LanguageCode: types.LanguageCode(c.opts.LanguageCode),
	})
	if err != nil {
		return nil, err
	}

	if len(c.opts.Labels) == 0 {
		for _, label := range output.Labels {
			if aws.ToFloat32(label.Score) >= c.opts.Threshold {
				return nil, errors.New("pii content found")
			}
		}
	} else {
		for _, label := range output.Labels {
			if util.Contains(c.opts.Labels, string(label.Name)) {
				if aws.ToFloat32(label.Score) >= c.opts.Threshold {
					return nil, errors.New("pii content found")
				}
			}
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: text,
	}, nil
}

func (c *AmazonComprehendPII) detectPII(ctx context.Context, text string) (schema.ChainValues, error) {
	output, err := c.client.DetectPiiEntities(ctx, &comprehend.DetectPiiEntitiesInput{
		Text:         aws.String(text),
		LanguageCode: types.LanguageCode(c.opts.LanguageCode),
	})
	if err != nil {
		return nil, err
	}

	if len(c.opts.Labels) == 0 {
		for _, entity := range output.Entities {
			if aws.ToFloat32(entity.Score) >= c.opts.Threshold {
				text = c.opts.RedactFunc(ctx, text, c.opts.MaskMarker, string(entity.Type), aws.ToInt32(entity.BeginOffset), aws.ToInt32(entity.EndOffset))
			}
		}
	} else {
		for _, entity := range output.Entities {
			if util.Contains(c.opts.Labels, string(entity.Type)) {
				if aws.ToFloat32(entity.Score) >= c.opts.Threshold {
					text = c.opts.RedactFunc(ctx, text, c.opts.MaskMarker, string(entity.Type), aws.ToInt32(entity.BeginOffset), aws.ToInt32(entity.EndOffset))
				}
			}
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: text,
	}, nil
}
