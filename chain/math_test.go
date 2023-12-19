package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestMath(t *testing.T) {
	t.Run("Valid Question", func(t *testing.T) {
		fake := llm.NewSimpleFake("```text\n3 * 3\n```")

		mathChain, err := NewMath(fake)
		assert.NoError(t, err)

		output, err := golc.SimpleCall(context.Background(), mathChain, "What is 3 times 3?")
		assert.NoError(t, err)
		assert.Equal(t, "9", output)
	})

	t.Run("Invalid Input Key", func(t *testing.T) {
		fake := llm.NewSimpleFake("```text\n3 * 3\n```")

		mathChain, err := NewMath(fake)
		assert.NoError(t, err)

		_, err = golc.Call(context.Background(), mathChain, schema.ChainValues{"invalid_key": "foo"})
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid chain values: no value for key question")
	})
}
