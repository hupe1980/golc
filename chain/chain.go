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

func Apply(ctx context.Context, chain schema.Chain, inputs []schema.ChainValues) ([]schema.ChainValues, error) {
	chainValues := []schema.ChainValues{}

	for _, input := range inputs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			vals, err := chain.Call(ctx, input)
			if err != nil {
				return nil, err
			}

			chainValues = append(chainValues, vals)
		}
	}

	return chainValues, nil
}
