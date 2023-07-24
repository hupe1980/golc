package outputparser

import (
	"testing"

	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestCommaSeparatedList(t *testing.T) {
	parser := NewCommaSeparatedList()

	// Test ParseResult
	t.Run("ParseResult", func(t *testing.T) {
		// Test case with a valid comma-separated list.
		result := schema.Generation{Text: "foo, bar, baz"}
		expected := []string{"foo", "bar", "baz"}
		actual, err := parser.ParseResult(result)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		// Test case with an empty result.
		result = schema.Generation{Text: ""}
		_, err = parser.ParseResult(result)
		assert.Error(t, err)
	})

	// Test Parse
	t.Run("Parse", func(t *testing.T) {
		// Test case with a valid comma-separated list.
		text := "foo, bar, baz"
		expected := []string{"foo", "bar", "baz"}
		actual, err := parser.Parse(text)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		// Test case with an empty text.
		text = ""
		_, err = parser.Parse(text)
		assert.Error(t, err)
	})

	// Test ParseWithPrompt
	t.Run("ParseWithPrompt", func(t *testing.T) {
		// Test case with a valid comma-separated list.
		text := "foo, bar, baz"
		p := prompt.StringPromptValue("dummy")
		expected := []string{"foo", "bar", "baz"}
		actual, err := parser.ParseWithPrompt(text, p)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		// Test case with an empty text.
		text = ""
		p = prompt.StringPromptValue("dummy")
		_, err = parser.ParseWithPrompt(text, p)
		assert.Error(t, err)
	})

	// Test GetFormatInstructions
	t.Run("GetFormatInstructions", func(t *testing.T) {
		expected := "Your response should be a list of comma-separated values, e.g.: `foo, bar, baz`"
		actual := parser.GetFormatInstructions()
		assert.Equal(t, expected, actual)
	})

	// Test Type
	t.Run("Type", func(t *testing.T) {
		expected := "comma_separated_list"
		actual := parser.Type()
		assert.Equal(t, expected, actual)
	})
}
