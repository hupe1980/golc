package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	t.Run("Format - no partial", func(t *testing.T) {
		prompt, err := NewTemplate("XX {{.foo}} XX")
		assert.NoError(t, err)

		text, err := prompt.Format(map[string]any{
			"foo": "foo",
		})
		assert.NoError(t, err)

		assert.Equal(t, "XX foo XX", text)
	})

	t.Run("Format - using partial", func(t *testing.T) {
		prompt, err := NewTemplate("XX {{.foo}}{{.bar}} XX", func(o *TemplateOptions) {
			o.PartialValues = PartialValues{"bar": "bar"}
		})
		assert.NoError(t, err)

		text, err := prompt.Format(map[string]any{
			"foo": "foo",
		})
		assert.NoError(t, err)

		assert.Equal(t, "XX foobar XX", text)
	})

	t.Run("Format - using partial func", func(t *testing.T) {
		prompt, err := NewTemplate("XX {{.foo}}{{.bar}} XX")
		assert.NoError(t, err)

		partialPrompt, err := prompt.Partial(PartialValues{"bar": "bar"})
		assert.NoError(t, err)

		text, err := partialPrompt.Format(map[string]any{
			"foo": "foo",
		})
		assert.NoError(t, err)

		assert.Equal(t, "XX foobar XX", text)
	})

	t.Run("Format - using full partial", func(t *testing.T) {
		prompt, err := NewTemplate("XX {{.foo}}{{.bar}} XX", func(o *TemplateOptions) {
			o.PartialValues = PartialValues{"foo": "foo", "bar": "bar"}
		})
		assert.NoError(t, err)

		text, err := prompt.Format(nil)
		assert.NoError(t, err)

		assert.Equal(t, "XX foobar XX", text)
	})

	t.Run("FormatPrompt", func(t *testing.T) {
		prompt, err := NewTemplate("XX {{.foo}} XX")
		assert.NoError(t, err)

		pv, err := prompt.FormatPrompt(map[string]any{
			"foo": "foo",
		})
		assert.NoError(t, err)

		assert.Equal(t, "XX foo XX", pv.String())
	})
}
