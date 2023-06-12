package prompt

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	pt, err := NewFormatter("This is a {{ .foo }} test.")
	assert.NoError(t, err)

	result, err := pt.Render(map[string]any{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, result, "This is a bar test.")
}

func TestListTemplateFields(t *testing.T) {
	template := template.Must(template.New("template").Parse("This is a {{ .foo }} test."))
	assert.ElementsMatch(t, ListTemplateFields(template), []string{"{{.foo}}"})
}
