package outputparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFencedCodeBlock_Parse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		// Create a test FencedCodeBlock instance
		fencedCodeBlock := NewFencedCodeBlock("```go")

		t.Run("Parse returns code blocks from text", func(t *testing.T) {
			// Define test input
			text := "Some text\n\n```go\nline1\nline2\n```\n\nMore text\n\n```python\nline3\nline4\n```\n\nFinal text"

			// Call the Parse method
			codeBlocks, err := fencedCodeBlock.Parse(text)
			assert.NoError(t, err)

			// Assert that the code blocks are extracted correctly
			expectedCodeBlocks := []string{"line1", "line2"}
			assert.ElementsMatch(t, expectedCodeBlocks, codeBlocks)
		})

		t.Run("Parse returns an error if fence is not found", func(t *testing.T) {
			// Define test input without the fence
			text := "Some text\n\n```typescript\nline1\nline2\n```\n\nMore text\n\n```python\nline3\nline4\n```\n\nFinal text"

			// Call the Parse method
			_, err := fencedCodeBlock.Parse(text)

			// Assert that an error is returned
			assert.Error(t, err)
		})
	})

	t.Run("GetFormatInstructions", func(t *testing.T) {
		// Create a test FencedCodeBlock instance
		fencedCodeBlock := NewFencedCodeBlock("```")

		t.Run("GetFormatInstructions returns an error", func(t *testing.T) {
			// Call the GetFormatInstructions method
			_, err := fencedCodeBlock.GetFormatInstructions()

			// Assert that an error is returned
			assert.Error(t, err)
		})
	})

	t.Run("Type", func(t *testing.T) {
		// Create a test FencedCodeBlock instance
		fencedCodeBlock := NewFencedCodeBlock("```")

		t.Run("Type returns the correct type", func(t *testing.T) {
			// Call the Type method
			typ := fencedCodeBlock.Type()

			// Assert that the correct type is returned
			assert.Equal(t, "fenced-code-block-output-parser", typ)
		})
	})
}
