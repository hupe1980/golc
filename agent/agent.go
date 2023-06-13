package agent

import (
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
)

type AgentType string

const (
	ZeroShotReactDescriptionAgentType AgentType = "zero-shot-react-description"
)

func toolNames(tools []golc.Tool) string {
	toolNames := []string{}
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name())
	}

	return strings.Join(toolNames, ", ")
}

func toolDescriptions(tools []golc.Tool) string {
	toolDescriptions := []string{}
	for _, tool := range tools {
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
	}

	return strings.Join(toolDescriptions, "\n")
}
