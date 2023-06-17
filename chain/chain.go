package chain

import (
	"context"
	"strings"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
)

type callbackOptions struct {
	Callbacks []schema.Callback
	Verbose   bool
}

type callFunc func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error)

type baseChain struct {
	chainName       string
	callFunc        callFunc
	inputKeys       []string
	outputKeys      []string
	memory          schema.Memory
	callbackOptions *callbackOptions
}

func (bc *baseChain) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	cm := callback.NewManager(bc.callbackOptions.Callbacks, bc.callbackOptions.Verbose)

	if err := cm.OnChainStart(bc.chainName, &inputs); err != nil {
		return nil, err
	}

	output, err := bc.callFunc(ctx, inputs)
	if err != nil {
		if cbError := cm.OnChainError(err); cbError != nil {
			return nil, cbError
		}

		return nil, err
	}

	if err := cm.OnChainEnd(&output); err != nil {
		return nil, err
	}

	return output, nil
}

func (bc *baseChain) Run(ctx context.Context, input any) (string, error) {
	if len(bc.inputKeys) != 1 {
		return "", ErrMultipleInputsInRun
	}

	if len(bc.outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{bc.inputKeys[0]: input}

	// TODO
	if bc.memory != nil {
		_, _ = bc.memory.LoadMemoryVariables(inputValues)
	}

	outputValues, err := bc.Call(ctx, inputValues)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[bc.outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return strings.TrimSpace(outputValue), nil
}

func (bc *baseChain) Apply(ctx context.Context, inputs []schema.ChainValues) ([]schema.ChainValues, error) {
	chainValues := []schema.ChainValues{}

	for _, input := range inputs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			vals, err := bc.Call(ctx, input)
			if err != nil {
				return nil, err
			}

			chainValues = append(chainValues, vals)
		}
	}

	return chainValues, nil
}

// InputKeys returns the expected input keys.
func (bc *baseChain) InputKeys() []string {
	return bc.inputKeys
}

// OutputKeys returns the output keys the chain will return.
func (bc *baseChain) OutputKeys() []string {
	return bc.outputKeys
}
