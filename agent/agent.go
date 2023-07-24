// Package agent provides functionality for creating and managing agents
// that leverage Large Language Models (LLMs) to make informed decisions and take actions.
package agent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hupe1980/golc/schema"
)

// toolNames returns a comma-separated string containing the names of the tools
// in the provided slice of schema.Tool.
func toolNames(tools []schema.Tool) string {
	toolNames := []string{}
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name())
	}

	return strings.Join(toolNames, ", ")
}

// toolDescriptions returns a formatted string containing the names and descriptions
// of the tools in the provided slice of schema.Tool. Each tool's name and description
// are listed in bullet points.
func toolDescriptions(tools []schema.Tool) string {
	toolDescriptions := []string{}
	for _, tool := range tools {
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
	}

	return strings.Join(toolDescriptions, "\n")
}

// inputsToString converts the values of the input map to strings and returns a new map with string values.
// The function takes a map of mixed data types (any) and converts the values to strings based on their types.
// Supported data types for conversion: string, int, int64, float32, float64, bool.
// For unknown data types, the function returns an error with ErrInputNotString.
// The returned map contains the keys from the input map with their corresponding string values.
// If any value in the input map cannot be converted to a string, the function returns an error with ErrInputNotString.
func inputsToString(inputValues map[string]any) (map[string]string, error) {
	inputs := make(map[string]string, len(inputValues))

	for key, value := range inputValues {
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int:
			valueStr = strconv.Itoa(v)
		case int64:
			valueStr = strconv.FormatInt(v, 10)
		case float32:
			valueStr = strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			valueStr = strconv.FormatBool(v)
		default:
			return nil, fmt.Errorf("%w: %s", ErrInputNotString, key)
		}

		inputs[key] = valueStr
	}

	return inputs, nil
}
