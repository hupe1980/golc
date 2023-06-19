package chain

import (
	"context"

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

	if bc.memory != nil {
		vars, _ := bc.memory.LoadMemoryVariables(ctx, inputs)
		for k, v := range vars {
			inputs[k] = v
		}
	}

	outputs, err := bc.callFunc(ctx, inputs)
	if err != nil {
		if cbError := cm.OnChainError(err); cbError != nil {
			return nil, cbError
		}

		return nil, err
	}

	if bc.memory != nil {
		if err := bc.memory.SaveContext(ctx, inputs, outputs); err != nil {
			return nil, err
		}
	}

	if err := cm.OnChainEnd(&outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

func (bc *baseChain) Run(ctx context.Context, input any) (string, error) {
	if len(bc.inputKeys) != 1 {
		return "", ErrMultipleInputsInRun
	}

	if len(bc.outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	outputValues, err := bc.Call(ctx, map[string]any{bc.inputKeys[0]: input})
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[bc.outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return outputValue, nil
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

func Call(ctx context.Context, chain schema.Chain, inputs schema.ChainValues) (schema.ChainValues, error) {
	cm := callback.NewManager(chain.Callbacks(), chain.Verbose())

	if err := cm.OnChainStart(chain.Type(), &inputs); err != nil {
		return nil, err
	}

	if chain.Memory() != nil {
		vars, _ := chain.Memory().LoadMemoryVariables(ctx, inputs)
		for k, v := range vars {
			inputs[k] = v
		}
	}

	outputs, err := chain.Call(ctx, inputs)
	if err != nil {
		if cbError := cm.OnChainError(err); cbError != nil {
			return nil, cbError
		}

		return nil, err
	}

	if chain.Memory() != nil {
		if err := chain.Memory().SaveContext(ctx, inputs, outputs); err != nil {
			return nil, err
		}
	}

	if err := cm.OnChainEnd(&outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

func Run(ctx context.Context, chain schema.Chain, input any) (string, error) {
	if len(chain.InputKeys()) != 1 {
		return "", ErrMultipleInputsInRun
	}

	if len(chain.OutputKeys()) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	outputValues, err := Call(ctx, chain, map[string]any{chain.InputKeys()[0]: input})
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[chain.OutputKeys()[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}

	return outputValue, nil
}
