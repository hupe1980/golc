package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate2(t *testing.T) {
	t.Run("Format", func(t *testing.T) {
		templateString := "Hello, {{.name}}! Your age is {{.age}}."
		template, err := NewTemplate(templateString)
		assert.NoError(t, err, "NewTemplate should not return an error")

		t.Run("Success", func(t *testing.T) {
			values := map[string]interface{}{
				"name": "John",
				"age":  30,
			}
			expectedResult := "Hello, John! Your age is 30."

			result, err := template.Format(values)
			assert.NoError(t, err, "Format should not return an error")
			assert.Equal(t, expectedResult, result, "Formatted result should match the expected result")
		})

		t.Run("WithPartialValues", func(t *testing.T) {
			partialValues := PartialValues{
				"name": "Jane",
			}
			updatedTemplate, err := template.Partial(partialValues)
			assert.NoError(t, err, "Partial should not return an error")

			values := map[string]interface{}{
				"age": 25,
			}
			expectedResult := "Hello, Jane! Your age is 25."

			result, err := updatedTemplate.Format(values)
			assert.NoError(t, err, "Format should not return an error")
			assert.Equal(t, expectedResult, result, "Formatted result should match the expected result with partial values")
		})
	})

	t.Run("InputVariables", func(t *testing.T) {
		templateString := "Hello, {{.name}}! Your age is {{.age}}."
		template, err := NewTemplate(templateString)
		assert.NoError(t, err, "NewTemplate should not return an error")

		t.Run("Success", func(t *testing.T) {
			expectedVariables := []string{"name", "age"}

			variables := template.InputVariables()
			assert.Equal(t, expectedVariables, variables, "Input variables should match the expected variables")
		})
	})

	t.Run("FormatPrompt", func(t *testing.T) {
		templateString := "What is your name? ({{.name}})"
		template, err := NewTemplate(templateString)
		assert.NoError(t, err, "NewTemplate should not return an error")

		t.Run("Success", func(t *testing.T) {
			values := map[string]interface{}{
				"name": "John",
			}
			expectedPrompt := "What is your name? (John)"

			promptValue, err := template.FormatPrompt(values)
			assert.NoError(t, err, "FormatPrompt should not return an error")
			assert.Equal(t, expectedPrompt, promptValue.String(), "Formatted prompt should match the expected prompt")
		})
	})
}
