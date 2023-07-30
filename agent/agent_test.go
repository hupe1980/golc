package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestToolNames(t *testing.T) {
	tools := []schema.Tool{
		&mockTool{ToolName: "Tool1"},
		&mockTool{ToolName: "Tool2"},
	}

	expected := "Tool1, Tool2"
	result := toolNames(tools)

	assert.Equal(t, expected, result)
}

func TestToolDescriptions(t *testing.T) {
	tools := []schema.Tool{
		&mockTool{ToolName: "Tool1", ToolDescription: "Description1."},
		&mockTool{ToolName: "Tool2", ToolDescription: "Description2."},
	}

	expected := `- Tool1: Description1.
- Tool2: Description2.`

	result := toolDescriptions(tools)

	assert.Equal(t, expected, result)
}

func TestInputsToString(t *testing.T) {
	inputValues := map[string]any{
		"param1": "value1",
		"param2": "value2",
	}

	expected := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	result, err := inputsToString(inputValues)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test case with non-string input value
	inputValues["param3"] = nil
	_, err = inputsToString(inputValues)
	assert.Error(t, err)
}

// Compile time check to ensure mockTool satisfies the Tool interface.
var _ schema.Tool = (*mockTool)(nil)

type mockTool struct {
	ToolName        string
	ToolDescription string
	ToolArgsType    any
	ToolRunFunc     func(ctx context.Context, input any) (string, error)
}

// Name returns the name of the tool.
func (t *mockTool) Name() string {
	if t.ToolName != "" {
		return t.ToolName
	}

	return "Mock"
}

// Description returns the description of the tool.
func (t *mockTool) Description() string {
	if t.ToolDescription != "" {
		return t.ToolDescription
	}

	return "Mock"
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *mockTool) ArgsType() reflect.Type {
	if t.ToolArgsType != nil {
		return reflect.TypeOf(t.ToolArgsType)
	}

	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *mockTool) Run(ctx context.Context, input any) (string, error) {
	if t.ToolRunFunc != nil {
		return t.ToolRunFunc(ctx, input)
	}

	return "Mock", nil
}

// Verbose returns the verbosity setting of the tool.
func (t *mockTool) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *mockTool) Callbacks() []schema.Callback {
	return nil
}
