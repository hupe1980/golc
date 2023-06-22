package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestSequential(t *testing.T) {
	t.Run("Call", func(t *testing.T) {
		ctx := context.Background()
		inputs := schema.ChainValues{
			"in1": "value1",
			"in2": "value2",
		}

		chain1 := &MockChain{
			CallFunc: func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
				return schema.ChainValues{
					"out1": "value1",
					"out2": "value2",
				}, nil
			},
			InputKeysFunc: func() []string {
				return []string{"in1", "in2"}
			},
			OutputKeysFunc: func() []string {
				return []string{"out1", "out2"}
			},
		}
		chain2 := &MockChain{
			CallFunc: func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
				return schema.ChainValues{
					"out3": "value3",
					"out4": "value4",
				}, nil
			},
			InputKeysFunc: func() []string {
				return []string{"out1", "out2"}
			},
			OutputKeysFunc: func() []string {
				return []string{"out3", "out4"}
			},
		}

		sequential, err := NewSequential([]schema.Chain{chain1, chain2}, []string{"in1", "in2"})
		assert.NoError(t, err)

		outputs, err := sequential.Call(ctx, inputs)
		assert.NoError(t, err)

		expectedOutputs := schema.ChainValues{
			"out3": "value3",
			"out4": "value4",
		}

		assert.Equal(t, expectedOutputs, outputs)
	})
}
