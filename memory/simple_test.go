package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	t.Run("MemoryKeys", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			memory := NewSimple()

			keys := memory.MemoryKeys()

			assert.Empty(t, keys, "Memory keys should be empty when no memories are stored")
		})

		t.Run("NonEmpty", func(t *testing.T) {
			memory := NewSimple()
			memory.memories["name"] = "John"
			memory.memories["age"] = 30

			expectedKeys := []string{"name", "age"}

			keys := memory.MemoryKeys()

			assert.ElementsMatch(t, expectedKeys, keys, "Memory keys should match the expected variables")
		})
	})

	t.Run("LoadMemoryVariables", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			memory := NewSimple()
			memory.memories["name"] = "John"
			memory.memories["age"] = 30

			inputs := map[string]interface{}{
				"var1": "value1",
				"var2": "value2",
			}

			expectedOutputs := map[string]interface{}{
				"name": "John",
				"age":  30,
			}

			outputs, err := memory.LoadMemoryVariables(context.TODO(), inputs)

			assert.NoError(t, err, "LoadMemoryVariables should not return an error")
			assert.Equal(t, expectedOutputs, outputs, "Loaded memory variables should match the expected outputs")
		})
	})

	t.Run("SaveContext", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			memory := NewSimple()

			inputs := map[string]interface{}{
				"name": "John",
				"age":  30,
			}

			outputs := map[string]interface{}{
				"var1": "value1",
				"var2": "value2",
			}

			err := memory.SaveContext(context.TODO(), inputs, outputs)

			assert.NoError(t, err, "SaveContext should not return an error")
			// No assertions made as SaveContext does not modify the state of the Simple memory.
		})
	})

	t.Run("Clear", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			memory := NewSimple()

			err := memory.Clear(context.TODO())

			assert.NoError(t, err, "Clear should not return an error")
			// No assertions made as Clear does not modify the state of the Simple memory.
		})
	})
}
