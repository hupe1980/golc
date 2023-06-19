package chain

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Transform satisfies the Chain interface.
var _ schema.Chain = (*TransformChain)(nil)

type TransformFunc func(inputs schema.ChainValues) (schema.ChainValues, error)

type TransformChainOptions struct {
	*schema.CallbackOptions
}

type TransformChain struct {
	inputKeys  []string
	outputKeys []string
	transform  TransformFunc
	opts       TransformChainOptions
}

func NewTransformChain(inputKeys, outputKeys []string, transform TransformFunc, optFns ...func(o *TransformChainOptions)) (*TransformChain, error) {
	opts := TransformChainOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &TransformChain{
		inputKeys:  inputKeys,
		outputKeys: outputKeys,
		transform:  transform,
		opts:       opts,
	}, nil
}

func (c *TransformChain) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return c.transform(inputs)
}

func (c *TransformChain) Memory() schema.Memory {
	return nil
}

func (c *TransformChain) Type() string {
	return "Transform"
}

func (c *TransformChain) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

func (c *TransformChain) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *TransformChain) InputKeys() []string {
	return c.inputKeys
}

// OutputKeys returns the output keys the chain will return.
func (c *TransformChain) OutputKeys() []string {
	return c.outputKeys
}
