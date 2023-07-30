// Package agent provides functionality for creating and managing agents
// that leverage Large Language Models (LLMs) to make informed decisions and take actions.
package agent

import (
	"fmt"
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
