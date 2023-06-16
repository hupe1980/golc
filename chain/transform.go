package chain

import (
	"context"

	"github.com/hupe1980/golc"
)

// Compile time check to ensure Transform satisfies the chain interface.
var _ golc.Chain = (*TransformChain)(nil)

type TransformFunc func(inputs golc.ChainValues) (golc.ChainValues, error)

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

func (t *TransformChain) Call(ctx context.Context, inputs golc.ChainValues) (golc.ChainValues, error) {
	return t.transform(inputs)
}
