package chain

import (
	"context"

	"github.com/hupe1980/golc"
)

type CallFunc func(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error)

type Chain struct {
	callFunc CallFunc
}

func NewChain(callFunc CallFunc) *Chain {
	return &Chain{
		callFunc: callFunc,
	}
}

func (b *Chain) Run(ctx context.Context) {}

func (b *Chain) Call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	return b.callFunc(ctx, values)
}

func (b *Chain) Apply(ctx context.Context, inputs []golc.ChainValues) ([]golc.ChainValues, error) {
	chainValues := []golc.ChainValues{}

	for _, input := range inputs {
		vals, err := b.Call(ctx, input)
		if err != nil {
			return nil, err
		}

		chainValues = append(chainValues, vals)
	}

	return chainValues, nil
}
