package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFewShotTemplate(t *testing.T) {
	// Define sample template, examples, and exampleTemplate
	template := "{{.Greeting}}, {{.Name}}!"
	examples := []map[string]interface{}{
		{"Greeting": "Hello"},
		{"Greeting": "Hi"},
	}
	exampleTemplate := NewTemplate("{{.Greeting}}")

	// Create a FewShotTemplate
	fsTemplate := NewFewShotTemplate(template, examples, exampleTemplate)

	t.Run("Format", func(t *testing.T) {
		values := map[string]interface{}{"Greeting": "Hey", "Name": "Charlie"}
		formatted, err := fsTemplate.Format(values)
		assert.NoError(t, err)
		assert.Equal(t, "Hello\n\nHi\n\nHey, Charlie!", formatted)
	})

	t.Run("FormatPrompt", func(t *testing.T) {
		values := map[string]interface{}{"Greeting": "Hey", "Name": "Charlie"}
		promptValue, err := fsTemplate.FormatPrompt(values)
		assert.NoError(t, err)
		assert.IsType(t, StringPromptValue(""), promptValue)
		assert.Equal(t, "Hello\n\nHi\n\nHey, Charlie!", promptValue.String())
	})

	t.Run("Partial", func(t *testing.T) {
		values := map[string]interface{}{"Greeting": "Hey"}
		partialValues := map[string]interface{}{"Name": "David"}
		partialTemplate := fsTemplate.Partial(partialValues)
		partialValues["Greeting"] = "Hi"
		formattedPartial, err := partialTemplate.Format(values)
		assert.NoError(t, err)
		assert.Equal(t, "Hello\n\nHi\n\nHey, David!", formattedPartial)
	})

	t.Run("OutputParser", func(t *testing.T) {
		outputParser, hasParser := fsTemplate.OutputParser()
		assert.False(t, hasParser)
		assert.Nil(t, outputParser)
	})

	t.Run("InputVariables", func(t *testing.T) {
		inputVars := fsTemplate.InputVariables()
		assert.ElementsMatch(t, inputVars, []string{"Greeting", "Name"})
	})
}
