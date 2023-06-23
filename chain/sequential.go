package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure Sequential satisfies the Chain interface.
var _ schema.Chain = (*Sequential)(nil)

type SequentialOptions struct {
	*schema.CallbackOptions
	Memory     schema.Memory
	OutputKeys []string
	ReturnAll  bool
}

type Sequential struct {
	chains     []schema.Chain
	inputKeys  []string
	outputKeys []string
	opts       SequentialOptions
}

func NewSequential(chains []schema.Chain, inputKeys []string, optFns ...func(o *SequentialOptions)) (*Sequential, error) {
	opts := SequentialOptions{
		ReturnAll: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	memoryKeys := []string{}
	if opts.Memory != nil {
		memoryKeys = opts.Memory.MemoryKeys()

		overlap := util.Intersect(inputKeys, memoryKeys)
		if len(overlap) > 0 {
			return nil, fmt.Errorf("overlapping input keys: %s", strings.Join(overlap, ","))
		}
	}

	knownKeys := append(inputKeys, memoryKeys...)

	for _, chain := range chains {
		missingKeys, _ := util.Difference(chain.InputKeys(), knownKeys)
		if len(missingKeys) > 0 {
			return nil, fmt.Errorf("missing required input keys: %s", strings.Join(missingKeys, ","))
		}

		overlap := util.Intersect(knownKeys, chain.OutputKeys())
		if len(overlap) > 0 {
			return nil, fmt.Errorf("overlapping output keys: %s", strings.Join(overlap, ","))
		}

		knownKeys = append(knownKeys, chain.OutputKeys()...)
	}

	if len(opts.OutputKeys) == 0 {
		if opts.ReturnAll {
			opts.OutputKeys, _ = util.Difference(knownKeys, inputKeys)
		} else {
			opts.OutputKeys = chains[len(chains)-1].OutputKeys()
		}
	}

	return &Sequential{
		chains:     chains,
		inputKeys:  inputKeys,
		outputKeys: opts.OutputKeys,
		opts:       opts,
	}, nil
}

// Call executes the Sequential chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *Sequential) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	knownValues := util.CopyMap(inputs)

	for _, c := range c.chains {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			outputs, err := golc.Call(ctx, c, knownValues)
			if err != nil {
				return nil, err
			}

			for k, v := range outputs {
				knownValues[k] = v
			}
		}
	}

	result := make(schema.ChainValues)
	for _, k := range c.opts.OutputKeys {
		result[k] = knownValues[k]
	}

	return result, nil
}

// Memory returns the memory associated with the chain.
func (c *Sequential) Memory() schema.Memory {
	return c.opts.Memory
}

// Type returns the type of the chain.
func (c *Sequential) Type() string {
	return "Sequential"
}

// Verbose returns the verbosity setting of the chain.
func (c *Sequential) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *Sequential) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Sequential) InputKeys() []string {
	return c.inputKeys
}

// OutputKeys returns the output keys the chain will return.
func (c *Sequential) OutputKeys() []string {
	return c.outputKeys
}

// Compile time check to ensure SimpleSequential satisfies the Chain interface.
var _ schema.Chain = (*SimpleSequential)(nil)

type SimpleSequentialOptions struct {
	*schema.CallbackOptions
	Memory       schema.Memory
	InputKey     string
	OutputKey    string
	StripOutputs bool
}

type SimpleSequential struct {
	chains []schema.Chain
	opts   SimpleSequentialOptions
}

func NewSimpleSequential(chains []schema.Chain, optFns ...func(o *SimpleSequentialOptions)) (*SimpleSequential, error) {
	opts := SimpleSequentialOptions{
		InputKey:     "input",
		OutputKey:    "output",
		StripOutputs: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	for _, chain := range chains {
		if len(chain.InputKeys()) != 1 {
			return nil, fmt.Errorf("chain with more than one expected input: %v", len(chain.InputKeys()))
		}

		if len(chain.OutputKeys()) != 1 {
			return nil, fmt.Errorf("chain with more than one expected output: %v", len(chain.OutputKeys()))
		}
	}

	return &SimpleSequential{
		chains: chains,
		opts:   opts,
	}, nil
}

// Call executes the SimpleSequential chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *SimpleSequential) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	input := inputs[c.opts.InputKey]

	for _, chain := range c.chains {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			input, err := golc.SimpleCall(ctx, chain, input)
			if err != nil {
				return nil, err
			}

			if c.opts.StripOutputs {
				input = strings.TrimSpace(input) //nolint ineffassign
			}
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: input,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *SimpleSequential) Memory() schema.Memory {
	return c.opts.Memory
}

// Type returns the type of the chain.
func (c *SimpleSequential) Type() string {
	return "SimpleSequential"
}

// Verbose returns the verbosity setting of the chain.
func (c *SimpleSequential) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *SimpleSequential) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *SimpleSequential) InputKeys() []string {
	return []string{}
}

// OutputKeys returns the output keys the chain will return.
func (c *SimpleSequential) OutputKeys() []string {
	return []string{}
}
