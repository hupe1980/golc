package agent

import (
	"fmt"
	"strings"

	"github.com/hupe1980/golc/schema"
)

func toolNames(tools []schema.Tool) string {
	toolNames := []string{}
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name())
	}

	return strings.Join(toolNames, ", ")
}

func toolDescriptions(tools []schema.Tool) string {
	toolDescriptions := []string{}
	for _, tool := range tools {
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
	}

	return strings.Join(toolDescriptions, "\n")
}

func inputsToString(inputValues map[string]any) (map[string]string, error) {
	inputs := make(map[string]string, len(inputValues))

	for key, value := range inputValues {
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrExecutorInputNotString, key)
		}

		inputs[key] = valueStr
	}

	return inputs, nil
}
