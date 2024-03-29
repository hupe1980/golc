package tool

import (
	"testing"

	"github.com/hupe1980/golc/integration/jsonschema"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestToOpenAIFunction(t *testing.T) {
	testCases := []struct {
		name          string
		inputTool     schema.Tool
		expectedFunc  *schema.FunctionDefinition
		expectedError error
	}{
		{
			name:      "Valid Tool",
			inputTool: &Sleep{},
			expectedFunc: &schema.FunctionDefinition{
				Name:        "Sleep",
				Description: "Make agent sleep for a specified number of seconds.",
				Parameters: schema.FunctionDefinitionParameters{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"__arg1": {
							Type:        "string",
							Description: "__arg1",
						},
					},
					Required: []string{"__arg1"},
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualFunc, err := ToFunction(tc.inputTool)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedFunc, actualFunc)
		})
	}
}
