package tool

import (
	"context"
	"errors"
	"reflect"

	huggingface "github.com/hupe1980/go-huggingface"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure HuggingFaceInjectionDetector satisfies the Tool interface.
var _ schema.Tool = (*HuggingFaceInjectionDetector)(nil)

// HuggingFaceInjectionDetectorClient is an interface for interacting with the Hugging Face injection detector.
type HuggingFaceInjectionDetectorClient interface {
	TextClassification(ctx context.Context, req *huggingface.TextClassificationRequest) (huggingface.TextClassificationResponse, error)
}

// HuggingFaceInjectionDetectorOptions represents configuration options for the Hugging Face injection detector.
type HuggingFaceInjectionDetectorOptions struct {
	// Model to use for injection detection.
	Model string
	// Threshold for injection detection.
	Threshold float32
	// Options represents optional settings for the classification.
	Options huggingface.Options
}

// HuggingFaceInjectionDetector represents a tool for detecting injection attacks using Hugging Face models.
type HuggingFaceInjectionDetector struct {
	client HuggingFaceInjectionDetectorClient
	opts   HuggingFaceInjectionDetectorOptions
}

// NewHuggingFaceInjectionDetector creates a new instance of the HuggingFaceInjectionDetector tool.
func NewHuggingFaceInjectionDetector(token string, optFns ...func(o *HuggingFaceInjectionDetectorOptions)) *HuggingFaceInjectionDetector {
	client := huggingface.NewInferenceClient(token)
	return NewHuggingFaceInjectionDetectorFromClient(client, optFns...)
}

// NewHuggingFaceInjectionDetectorFromClient creates a new instance of the HuggingFaceInjectionDetector tool.
func NewHuggingFaceInjectionDetectorFromClient(client HuggingFaceInjectionDetectorClient, optFns ...func(o *HuggingFaceInjectionDetectorOptions)) *HuggingFaceInjectionDetector {
	opts := HuggingFaceInjectionDetectorOptions{
		Model:     "deepset/deberta-v3-base-injection",
		Threshold: 0.8,
		Options:   huggingface.Options{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &HuggingFaceInjectionDetector{
		client: client,
		opts:   opts,
	}
}

// Name returns the name of the tool.
func (t *HuggingFaceInjectionDetector) Name() string {
	return "HuggingFaceInjectionDetector"
}

// Description returns the description of the tool.
func (t *HuggingFaceInjectionDetector) Description() string {
	return `A wrapper around HuggingFace Prompt Injection security model.
Useful for when you need to ensure that prompt is free of injection attacks.
Input should be any message from the user.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *HuggingFaceInjectionDetector) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *HuggingFaceInjectionDetector) Run(ctx context.Context, input any) (string, error) {
	query, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	resp, err := t.client.TextClassification(ctx, &huggingface.TextClassificationRequest{
		Inputs:  query,
		Model:   t.opts.Model,
		Options: t.opts.Options,
	})
	if err != nil {
		return "", err
	}

	for _, v := range resp[0] {
		if v.Label == "INJECTION" && v.Score >= t.opts.Threshold {
			return "", errors.New("prompt injection attack detected")
		}
	}

	return query, nil
}

// Verbose returns the verbosity setting of the tool.
func (t *HuggingFaceInjectionDetector) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *HuggingFaceInjectionDetector) Callbacks() []schema.Callback {
	return nil
}
