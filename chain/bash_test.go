package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBash(t *testing.T) {
	t.Run("Valid Question", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			return "```bash\necho 'hello world'\n```"
		})

		bashChain, err := NewBash(fake, func(o *BashOptions) {
			o.BashRunner = &mockBashRunner{
				Output: "hello world",
				Error:  nil,
			}
		})
		require.NoError(t, err)

		output, err := golc.SimpleCall(context.Background(), bashChain, "Please write a bash script that prints 'Hello World' to the console.")
		assert.NoError(t, err)
		assert.Contains(t, "hello world", output)
	})

	t.Run("Invalid Input Key", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			return "```bash\necho 'hello world'\n```"
		})

		bashChain, err := NewBash(fake, func(o *BashOptions) {
			o.BashRunner = &mockBashRunner{
				Output: "hello world",
				Error:  nil,
			}
		})
		require.NoError(t, err)

		_, err = golc.Call(context.Background(), bashChain, schema.ChainValues{"invalid_key": "foo"})
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid input values: no value for inputKey question")
	})

	t.Run("Invalid commands", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			return "```bash\necho 'hello world'\n```"
		})

		bashChain, err := NewBash(fake, func(o *BashOptions) {
			o.BashRunner = &mockBashRunner{
				Output: "hello world",
				Error:  nil,
			}
			o.VerifyCommands = func(commands []string) bool { return false }
		})
		require.NoError(t, err)

		_, err = golc.SimpleCall(context.Background(), bashChain, "Please write a bash script that prints 'Hello World' to the console.")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid commands: [echo 'hello world']")
	})
}

type mockBashRunner struct {
	Output string
	Error  error
}

func (mbr *mockBashRunner) Run(ctx context.Context, commands []string) (string, error) {
	if mbr.Error != nil {
		return "", mbr.Error
	}

	return mbr.Output, nil
}
