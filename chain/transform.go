package chain

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Transform satisfies the Chain interface.
var _ schema.Chain = (*TransformChain)(nil)

type TransformFunc func(inputs schema.ChainValues) (schema.ChainValues, error)

type TransformChain struct {
	*chain
	inputKeys  []string
	outputKeys []string
	transform  TransformFunc
}

func NewTransformChain(inputKeys, outputKeys []string, transform TransformFunc) (*TransformChain, error) {
	t := &TransformChain{
		inputKeys:  inputKeys,
		outputKeys: outputKeys,
		transform:  transform,
	}

	t.chain = newChain(t.call, t.inputKeys, t.outputKeys)

	return t, nil
}

func (t *TransformChain) call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return t.transform(inputs)
}
