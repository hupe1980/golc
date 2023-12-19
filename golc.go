// Package golc provides functions for executing chains.
package golc

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"golang.org/x/sync/errgroup"
)

var (
	// Verbose controls the verbosity of the chain execution.
	Verbose = false
)

type CallOptions struct {
	Callbacks      []schema.Callback
	ParentRunID    string
	IncludeRunInfo bool
	Stop           []string
}

// Call executes a chain with multiple inputs.
// It returns the outputs of the chain or an error, if any.
func Call(ctx context.Context, chain schema.Chain, inputs schema.ChainValues, optFns ...func(*CallOptions)) (schema.ChainValues, error) {
	opts := CallOptions{
		IncludeRunInfo: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	cm := callback.NewManager(opts.Callbacks, chain.Callbacks(), chain.Verbose(), func(mo *callback.ManagerOptions) {
		mo.ParentRunID = opts.ParentRunID
	})

	rm, err := cm.OnChainStart(ctx, &schema.ChainStartManagerInput{
		ChainType: chain.Type(),
		Inputs:    inputs,
	})
	if err != nil {
		return nil, err
	}

	if chain.Memory() != nil {
		vars, _ := chain.Memory().LoadMemoryVariables(ctx, inputs)
		for k, v := range vars {
			inputs[k] = v
		}
	}

	outputs, err := chain.Call(ctx, inputs, func(o *schema.CallOptions) {
		o.CallbackManger = rm
		o.Stop = opts.Stop
	})
	if err != nil {
		if cbErr := rm.OnChainError(ctx, &schema.ChainErrorManagerInput{
			Error: err,
		}); cbErr != nil {
			return nil, cbErr
		}

		return nil, err
	}

	if chain.Memory() != nil {
		if err := chain.Memory().SaveContext(ctx, inputs, outputs); err != nil {
			return nil, err
		}
	}

	if err := rm.OnChainEnd(ctx, &schema.ChainEndManagerInput{
		Outputs: outputs,
	}); err != nil {
		return nil, err
	}

	if opts.IncludeRunInfo {
		outputs["runInfo"] = cm.RunID()
	}

	return outputs, nil
}

type SimpleCallOptions struct {
	Callbacks   []schema.Callback
	ParentRunID string
	Stop        []string
}

// SimpleCall executes a chain with a single input and a single output.
// It returns the output value as a string or an error, if any.
func SimpleCall(ctx context.Context, chain schema.Chain, input any, optFns ...func(*SimpleCallOptions)) (string, error) {
	opts := SimpleCallOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	var cv schema.ChainValues

	switch v := input.(type) {
	case schema.ChainValues:
		cv = v
	default:
		if len(chain.InputKeys()) != 1 {
			return "", fmt.Errorf("invalid arguments: number of input keys must be 1, got %d", len(chain.InputKeys()))
		}

		cv = schema.ChainValues{
			chain.InputKeys()[0]: input,
		}
	}

	if len(chain.OutputKeys()) != 1 {
		return "", fmt.Errorf("invalid arguments: number of output keys must be 1, got %d", len(chain.OutputKeys()))
	}

	outputValues, err := Call(ctx, chain, cv, func(o *CallOptions) {
		o.Callbacks = opts.Callbacks
		o.ParentRunID = opts.ParentRunID
		o.Stop = opts.Stop
	})
	if err != nil {
		return "", err
	}

	return outputValues.GetString(chain.OutputKeys()[0])
}

type BatchCallOptions struct {
	Callbacks      []schema.Callback
	ParentRunID    string
	IncludeRunInfo bool
	Stop           []string
	MaxConcurrency int
}

// BatchCall executes multiple calls to the chain.Call function concurrently and collects
// the results in the same order as the inputs. It utilizes the errgroup package to manage
// the concurrent execution and handle any errors that may occur.
func BatchCall(ctx context.Context, chain schema.Chain, inputs []schema.ChainValues, optFns ...func(*BatchCallOptions)) ([]schema.ChainValues, error) {
	opts := BatchCallOptions{
		MaxConcurrency: 5,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	errs, errctx := errgroup.WithContext(ctx)

	errs.SetLimit(opts.MaxConcurrency)

	chainValues := make([]schema.ChainValues, len(inputs))

	for i, input := range inputs {
		i, input := i, input

		errs.Go(func() error {
			vals, err := Call(errctx, chain, input, func(o *CallOptions) {
				o.Callbacks = opts.Callbacks
				o.ParentRunID = opts.ParentRunID
				o.IncludeRunInfo = opts.IncludeRunInfo
				o.Stop = opts.Stop
			})
			if err != nil {
				return err
			}

			chainValues[i] = vals

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return chainValues, nil
}
