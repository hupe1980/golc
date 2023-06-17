package chain

import (
	"context"
	"strings"

	"github.com/hupe1980/golc/schema"
)

type callFunc func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error)

type chain struct {
	callFunc   callFunc
	inputKeys  []string
	outputKeys []string
}

func newChain(callFunc callFunc, inputKeys []string, outputKeys []string) *chain {
	return &chain{
		callFunc:   callFunc,
		inputKeys:  inputKeys,
		outputKeys: outputKeys,
	}
}

func (c *chain) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	return c.callFunc(ctx, inputs)
}

func (c *chain) Run(ctx context.Context, input any) (string, error) {
	if len(c.inputKeys) != 1 {
		return "", ErrMultipleInputsInRun
	}

	if len(c.outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{c.inputKeys[0]: input}

	outputValues, err := c.Call(ctx, inputValues)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[c.outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return strings.TrimSpace(outputValue), nil
}

func (c *chain) Apply(ctx context.Context, inputs []schema.ChainValues) ([]schema.ChainValues, error) {
	chainValues := []schema.ChainValues{}

	for _, input := range inputs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			vals, err := c.Call(ctx, input)
			if err != nil {
				return nil, err
			}

			chainValues = append(chainValues, vals)
		}
	}

	return chainValues, nil
}

// InputKeys returns the expected input keys.
func (c *chain) InputKeys() []string {
	return c.inputKeys
}

// OutputKeys returns the output keys the chain will return.
func (c *chain) OutputKeys() []string {
	return c.outputKeys
}

type callbackOptions struct {
	Callbacks []schema.Callback
	Verbose   bool
}
