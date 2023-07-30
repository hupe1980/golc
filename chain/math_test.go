package chain

import (
	"context"
	"fmt"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestMath(t *testing.T) {
	t.Run("Valid Question", func(t *testing.T) {
		fake := llm.NewSimpleFake("```text\n3 * 3\n```")

		mathChain, err := NewMath(fake)
		require.NoError(t, err)

		output, err := golc.SimpleCall(context.Background(), mathChain, "What is 3 times 3?")
		require.NoError(t, err)

		fmt.Println(output)
		require.Equal(t, "9", output)
	})

	t.Run("Invalid Input Key", func(t *testing.T) {
		fake := llm.NewSimpleFake("```text\n3 * 3\n```")

		mathChain, err := NewMath(fake)
		require.NoError(t, err)

		_, err = golc.Call(context.Background(), mathChain, schema.ChainValues{"invalid_key": "foo"})
		require.Error(t, err)
		require.EqualError(t, err, "invalid input values: no value for inputKey question")
	})
}
