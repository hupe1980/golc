package toolkit

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

// TestNewBrowser tests the creation of a new Browser object
func TestNewBrowser(t *testing.T) {
	browser, err := NewBrowser(nil)
	require.NoError(t, err)
	require.NotNil(t, browser)

	// Ensure that the Browser has the expected tools
	expectedToolNames := []string{
		"CurrentPage",
		"NavigateBrowser",
		"ExtractText",
	}
	tools := browser.Tools()
	require.Len(t, tools, len(expectedToolNames))

	for _, name := range expectedToolNames {
		assertToolExists(t, tools, name)
	}
}

// assertToolExists checks if a tool with the specified name exists in the list of tools.
func assertToolExists(t *testing.T, tools []schema.Tool, name string) {
	t.Helper()

	for _, tool := range tools {
		if tool.Name() == name {
			return // Found the tool
		}
	}

	require.Fail(t, "tool not found: "+name)
}
