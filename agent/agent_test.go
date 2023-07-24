package agent

import (
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
	"github.com/stretchr/testify/assert"
)

func TestToolNames(t *testing.T) {
	tools := []schema.Tool{
		tool.NewSleep(),
		tool.NewHuman(),
	}

	expected := "Sleep, Human"
	result := toolNames(tools)

	assert.Equal(t, expected, result)
}

func TestToolDescriptions(t *testing.T) {
	tools := []schema.Tool{
		tool.NewSleep(),
		tool.NewHuman(),
	}

	expected := `- Sleep: Make agent sleep for a specified number of seconds.
- Human: You can ask a human for guidance when you think you got stuck or you are not sure what to do next. The input should be a question for the human.`

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
