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
	*callbackOptions
}

type TransformChain struct {
	*baseChain
	transform TransformFunc
	opts      TransformChainOptions
}

func NewTransformChain(inputKeys, outputKeys []string, transform TransformFunc, optFns ...func(o *TransformChainOptions)) (*TransformChain, error) {
	opts := TransformChainOptions{
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	t := &TransformChain{
		transform: transform,
		opts:      opts,
	}

	t.baseChain = &baseChain{
		chainName:       "TransformChain",
		callFunc:        t.call,
		inputKeys:       inputKeys,
		outputKeys:      outputKeys,
		callbackOptions: opts.callbackOptions,
	}

	return t, nil
}

func (t *TransformChain) call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return t.transform(inputs)
}
