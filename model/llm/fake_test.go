package llm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFake(t *testing.T) {
	// Define a sample response function for the Fake LLM model
	responseFunc := func(prompt string) string {
		return "Generated text based on prompt: " + prompt
	}

	// Create a new instance of the Fake LLM model
	fake := NewFake(responseFunc)

	// Test the Generate method
	t.Run("Generate", func(t *testing.T) {
		prompt := "Hello, world!"
		expectedText := "Generated text based on prompt: Hello, world!"

		result, err := fake.Generate(context.Background(), prompt)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Generations, 1)
		assert.Equal(t, expectedText, result.Generations[0].Text)
	})

	// Test the Type method
	t.Run("Type", func(t *testing.T) {
		expectedType := "llm.Fake"
		assert.Equal(t, expectedType, fake.Type())
	})

	// Test the Verbose method
	t.Run("Verbose", func(t *testing.T) {
		assert.False(t, fake.Verbose())
	})

	// Test the Callbacks method
	t.Run("Callbacks", func(t *testing.T) {
		callbacks := fake.Callbacks()
		assert.Empty(t, callbacks)
	})

	// Test the InvocationParams method
	t.Run("InvocationParams", func(t *testing.T) {
		invocationParams := fake.InvocationParams()
		assert.NotNil(t, invocationParams)
	})
}
