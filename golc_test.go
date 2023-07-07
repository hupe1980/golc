package golc

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	// Define the inputs and expected outputs
	inputs := schema.ChainValues{
		"input": "test",
	}
	expectedOutputs := schema.ChainValues{
		"output": "result",
	}

	// Create a mock chain
	chain := mockChain{
		CallFunc: func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
			// Validate the inputs
			assert.Equal(t, "test", inputs["input"])

			// Return the expected outputs
			return expectedOutputs, nil
		},
	}

	// Call the chain
	outputs, err := Call(context.Background(), chain, inputs)
	assert.NoError(t, err)

	// Validate the outputs
	assert.Equal(t, expectedOutputs, outputs)
}

func TestSimpleCall(t *testing.T) {
	// Define the input and expected output
	input := "test"
	expectedOutput := "result"

	// Create a mock chain
	chain := mockChain{
		CallFunc: func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
			// Validate the input
			assert.Equal(t, "test", inputs["input"])

			// Return the expected output
			return schema.ChainValues{
				"output": expectedOutput,
			}, nil
		},
		InputKeysFunc: func() []string {
			return []string{"input"}
		},
		OutputKeysFunc: func() []string {
			return []string{"output"}
		},
	}

	// Call the chain
	output, err := SimpleCall(context.Background(), chain, input)
	assert.NoError(t, err)

	// Validate the output
	assert.Equal(t, expectedOutput, output)
}

func TestBatchCall(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name          string
		ctx           context.Context
		chain         schema.Chain
		inputs        []schema.ChainValues // Define your input values here
		expected      []schema.ChainValues // Define the expected output values here
		expectedError error
	}{
		{
			name: "Success",
			ctx:  context.TODO(),
			chain: mockChain{
				CallFunc: func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
					return inputs, nil
				},
			},
			inputs: []schema.ChainValues{
				{"foo1": "bar1"}, {"foo2": "bar2"},
			},
			expected: []schema.ChainValues{
				{"foo1": "bar1"}, {"foo2": "bar2"},
			},
			expectedError: nil,
		},
		{
			name: "Error",
			ctx:  context.TODO(),
			chain: mockChain{
				CallFunc: func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
					return nil, errors.New("error occurred during chain.Call")
				},
			},
			inputs: []schema.ChainValues{
				{"foo1": "bar1"}, {"foo2": "bar2"},
			},
			expected:      nil,
			expectedError: errors.New("error occurred during chain.Call"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := BatchCall(tc.ctx, tc.chain, tc.inputs)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

// mockChain is a mock implementation of the schema.Chain interface
type mockChain struct {
	CallFunc       func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error)
	TypeFunc       func() string
	VerboseFunc    func() bool
	CallbacksFunc  func() []schema.Callback
	MemoryFunc     func() schema.Memory
	InputKeysFunc  func() []string
	OutputKeysFunc func() []string
}

// Call is the mock implementation of the Call method
func (m mockChain) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, inputs)
	}

	return schema.ChainValues{}, nil
}

// Type is the mock implementation of the Type method
func (m mockChain) Type() string {
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}

	return "Mock"
}

// Verbose is the mock implementation of the Verbose method
func (m mockChain) Verbose() bool {
	if m.VerboseFunc != nil {
		return m.VerboseFunc()
	}

	return false
}

// Callbacks is the mock implementation of the Callbacks method
func (m mockChain) Callbacks() []schema.Callback {
	if m.CallbacksFunc != nil {
		return m.CallbacksFunc()
	}

	return nil
}

// Memory is the mock implementation of the Memory method
func (m mockChain) Memory() schema.Memory {
	if m.MemoryFunc != nil {
		return m.MemoryFunc()
	}

	return nil
}

// InputKeys is the mock implementation of the InputKeys method
func (m mockChain) InputKeys() []string {
	if m.InputKeysFunc != nil {
		return m.InputKeysFunc()
	}

	return nil
}

// OutputKeys is the mock implementation of the OutputKeys method
func (m mockChain) OutputKeys() []string {
	if m.OutputKeysFunc != nil {
		return m.OutputKeysFunc()
	}

	return nil
}
