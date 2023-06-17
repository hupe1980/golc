package chain

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Transform satisfies the Chain interface.
var _ schema.Chain = (*TransformChain)(nil)

type TransformFunc func(inputs schema.ChainValues) (schema.ChainValues, error)

type TransformChain struct {
	inputKeys  []string
	outputKeys []string
	transform  TransformFunc
}

func NewTransformChain(inputKeys, outputKeys []string, transform TransformFunc) (*TransformChain, error) {
	return &TransformChain{
		inputKeys:  inputKeys,
		outputKeys: outputKeys,
		transform:  transform,
	}, nil
}

func (t *TransformChain) InputKeys() []string {
	return t.inputKeys
}

func (t *TransformChain) OutputKeys() []string {
	return t.outputKeys
}

func (t *TransformChain) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return t.transform(inputs)
}
