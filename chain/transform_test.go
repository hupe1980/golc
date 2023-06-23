package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	t.Run("InputKeys", func(t *testing.T) {
		// Test case: InputKeys returns the expected input keys
		inputKeys := []string{"input1", "input2"}
		outputKeys := []string{"output1", "output2"}

		transform := func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
			return nil, nil
		}

		chain, err := NewTransform(inputKeys, outputKeys, transform)
		assert.NoError(t, err)

		expectedInputKeys := []string{"input1", "input2"}

		result := chain.InputKeys()
		assert.Equal(t, expectedInputKeys, result)
	})

	t.Run("OutputKeys", func(t *testing.T) {
		// Test case: OutputKeys returns the expected output keys
		inputKeys := []string{"input1", "input2"}
		outputKeys := []string{"output1", "output2"}

		transform := func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
			return nil, nil
		}

		chain, err := NewTransform(inputKeys, outputKeys, transform)
		assert.NoError(t, err)

		expectedOutputKeys := []string{"output1", "output2"}

		result := chain.OutputKeys()
		assert.Equal(t, expectedOutputKeys, result)
	})

	t.Run("Call", func(t *testing.T) {
		// Test case: Call invokes the transform function and returns the result
		inputKeys := []string{"input1", "input2"}
		outputKeys := []string{"output1", "output2"}

		transform := func(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
			result := make(schema.ChainValues)
			result["output1"] = inputs["input1"].(string) + "-transformed"
			result["output2"] = inputs["input2"].(int) * 2
			return result, nil
		}

		chain, err := NewTransform(inputKeys, outputKeys, transform)
		assert.NoError(t, err)

		inputs := make(schema.ChainValues)
		inputs["input1"] = "value1"
		inputs["input2"] = 5

		expectedResult := make(schema.ChainValues)
		expectedResult["output1"] = "value1-transformed"
		expectedResult["output2"] = 10

		result, err := chain.Call(context.Background(), inputs)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})
}
