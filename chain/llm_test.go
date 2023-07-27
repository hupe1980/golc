package chain

import (
	"context"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLLM(t *testing.T) {
	t.Run("Valid Question", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			return "This is a valid question."
		})

		llmChain, err := NewLLM(fake, prompt.NewTemplate("{{.input}}"))
		require.NoError(t, err)

		output, err := golc.SimpleCall(context.Background(), llmChain, "Please provide a valid question.")
		assert.NoError(t, err)
		assert.Equal(t, output, "This is a valid question.")
	})
}
