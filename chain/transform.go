package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Transform satisfies the Chain interface.
var _ schema.Chain = (*Transform)(nil)

type TransformFunc func(inputs schema.ChainValues) (schema.ChainValues, error)

type TransformOptions struct {
	*schema.CallbackOptions
}

type Transform struct {
	inputKeys  []string
	outputKeys []string
	transform  TransformFunc
	opts       TransformOptions
}

func NewTransform(inputKeys, outputKeys []string, transform TransformFunc, optFns ...func(o *TransformOptions)) (*Transform, error) {
	opts := TransformOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Transform{
		inputKeys:  inputKeys,
		outputKeys: outputKeys,
		transform:  transform,
		opts:       opts,
	}, nil
}

func (c *Transform) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return c.transform(inputs)
}

func (c *Transform) Memory() schema.Memory {
	return nil
}

func (c *Transform) Type() string {
	return "Transform"
}

func (c *Transform) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

func (c *Transform) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Transform) InputKeys() []string {
	return c.inputKeys
}

// OutputKeys returns the output keys the chain will return.
func (c *Transform) OutputKeys() []string {
	return c.outputKeys
}
