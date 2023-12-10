package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadonly(t *testing.T) {
	ctx := context.Background()

	mockMemory := NewConversationBuffer()
	err := mockMemory.SaveContext(ctx, map[string]any{"input": "foo"}, map[string]any{"output": "bar"})
	assert.NoError(t, err)

	readonly := NewReadonly(mockMemory)

	t.Run("MemoryKeys", func(t *testing.T) {
		assert.Equal(t, mockMemory.MemoryKeys(), readonly.MemoryKeys())
	})

	t.Run("LoadMemoryVariables", func(t *testing.T) {
		mockInputs := map[string]any{"key": "value"}
		result, err := readonly.LoadMemoryVariables(ctx, mockInputs)
		assert.NoError(t, err)
		assert.Equal(t, "Human: foo\nAI: bar", result["history"])
	})

	t.Run("SaveContext", func(t *testing.T) {
		err := readonly.SaveContext(ctx, map[string]any{"input": "Not saved"}, map[string]any{"output": "Not saved"})
		assert.NoError(t, err)
		result, err := readonly.LoadMemoryVariables(ctx, map[string]any{})
		assert.NoError(t, err)
		assert.Equal(t, "Human: foo\nAI: bar", result["history"])
	})

	t.Run("Clear", func(t *testing.T) {
		err := readonly.Clear(ctx)
		assert.NoError(t, err)
		result, err := readonly.LoadMemoryVariables(ctx, map[string]any{})
		assert.NoError(t, err)
		assert.Equal(t, "Human: foo\nAI: bar", result["history"])
	})
}
