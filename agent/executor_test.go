package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	t.Parallel()

	tool := &mockTool{
		ToolRunFunc: func(ctx context.Context, input interface{}) (string, error) {
			return "Observation", nil
		},
	}

	t.Run("Call_Success", func(t *testing.T) {
		t.Parallel()

		inputs := schema.ChainValues{
			"key1": "value1",
			"key2": 42,
		}

		expectedOutputs := schema.ChainValues{"outputKey": "outputValue"}

		agent := &mockAgent{
			PlanFunc: func(ctx context.Context, steps []schema.AgentStep, inputs schema.ChainValues) ([]*schema.AgentAction, *schema.AgentFinish, error) {
				return []*schema.AgentAction{}, &schema.AgentFinish{
					ReturnValues: expectedOutputs,
				}, nil
			},
		}

		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		outputs, err := executor.Call(context.Background(), inputs)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutputs, outputs)
	})

	t.Run("Call_Success", func(t *testing.T) {
		t.Parallel()

		inputs := schema.ChainValues{
			"key1": "value1",
			"key2": 42,
		}

		agent := &mockAgent{
			PlanFunc: func(ctx context.Context, steps []schema.AgentStep, inputs schema.ChainValues) ([]*schema.AgentAction, *schema.AgentFinish, error) {
				return nil, nil, errors.New("executor error")
			},
		}

		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		_, err = executor.Call(context.Background(), inputs)
		assert.ErrorContains(t, err, "executor error")
	})

	t.Run("InputKeys", func(t *testing.T) {
		agent := &mockAgent{
			IKeys: []string{"foo", "bar"},
		}
		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		keys := executor.InputKeys()
		assert.ElementsMatch(t, keys, []string{"foo", "bar"})
	})

	t.Run("OutputKeys", func(t *testing.T) {
		agent := &mockAgent{
			OKeys: []string{"foo", "bar"},
		}
		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		keys := executor.OutputKeys()
		assert.ElementsMatch(t, keys, []string{"foo", "bar"})
	})

	t.Run("Type", func(t *testing.T) {
		agent := &mockAgent{}
		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		typ := executor.Type()
		assert.Equal(t, "Executor", typ)
	})

	t.Run("Verbose", func(t *testing.T) {
		agent := &mockAgent{}
		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		verbose := executor.Verbose()

		assert.Equal(t, executor.opts.CallbackOptions.Verbose, verbose)
	})

	t.Run("Callbacks", func(t *testing.T) {
		agent := &mockAgent{}
		executor, err := NewExecutor(agent, []schema.Tool{tool})
		assert.NoError(t, err)

		callbacks := executor.Callbacks()

		assert.Equal(t, executor.opts.CallbackOptions.Callbacks, callbacks)
	})
}

// mockAgent is a custom mock for the schema.Agent interface.
type mockAgent struct {
	IKeys    []string
	OKeys    []string
	PlanFunc func(ctx context.Context, steps []schema.AgentStep, inputs schema.ChainValues) ([]*schema.AgentAction, *schema.AgentFinish, error)
}

// Plan is a method required by the schema.Agent interface.
func (m *mockAgent) Plan(ctx context.Context, steps []schema.AgentStep, inputs schema.ChainValues) ([]*schema.AgentAction, *schema.AgentFinish, error) {
	if m.PlanFunc != nil {
		return m.PlanFunc(ctx, steps, inputs)
	}

	panic("plan func not implemented")
}

// Name is a method required by the schema.Agent interface.
func (m *mockAgent) Name() string {
	return "mockAgent"
}

// Callbacks is a method required by the schema.Agent interface.
func (m *mockAgent) Callbacks() []schema.Callback {
	return nil
}

// InputKeys is a method required by the schema.Agent interface.
func (m *mockAgent) InputKeys() []string {
	return m.IKeys
}

// OutputKeys is a method required by the schema.Agent interface.
func (m *mockAgent) OutputKeys() []string {
	return m.OKeys
}
