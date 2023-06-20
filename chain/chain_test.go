package chain

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

// MockChain is a mock implementation of the schema.Chain interface
type MockChain struct {
	CallFunc       func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error)
	TypeFunc       func() string
	VerboseFunc    func() bool
	CallbacksFunc  func() []schema.Callback
	MemoryFunc     func() schema.Memory
	InputKeysFunc  func() []string
	OutputKeysFunc func() []string
}

// Call is the mock implementation of the Call method
func (m MockChain) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, inputs)
	}

	return schema.ChainValues{}, nil
}

// Type is the mock implementation of the Type method
func (m MockChain) Type() string {
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}

	return "Mock"
}

// Verbose is the mock implementation of the Verbose method
func (m MockChain) Verbose() bool {
	if m.VerboseFunc != nil {
		return m.VerboseFunc()
	}

	return false
}

// Callbacks is the mock implementation of the Callbacks method
func (m MockChain) Callbacks() []schema.Callback {
	if m.CallbacksFunc != nil {
		return m.CallbacksFunc()
	}

	return nil
}

// Memory is the mock implementation of the Memory method
func (m MockChain) Memory() schema.Memory {
	if m.MemoryFunc != nil {
		return m.MemoryFunc()
	}

	return nil
}

// InputKeys is the mock implementation of the InputKeys method
func (m MockChain) InputKeys() []string {
	if m.InputKeysFunc != nil {
		return m.InputKeysFunc()
	}

	return nil
}

// OutputKeys is the mock implementation of the OutputKeys method
func (m MockChain) OutputKeys() []string {
	if m.OutputKeysFunc != nil {
		return m.OutputKeysFunc()
	}

	return nil
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
			chain: MockChain{
				CallFunc: func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
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
			chain: MockChain{
				CallFunc: func(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
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
			result, err := golc.BatchCall(tc.ctx, tc.chain, tc.inputs)

			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
