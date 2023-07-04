package agent

import (
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

type AgentType string

const (
	ReactDescriptionAgentType               AgentType = "react-description"
	ReactDocstoreAgentType                  AgentType = "react-docstore"
	ConversationalReactDescriptionAgentType AgentType = "conversational-react-description"
	OpenAIFunctionsAgentType                AgentType = "openai-functions"
)

type Options struct {
	*schema.CallbackOptions
}

func New(llm schema.Model, tools []schema.Tool, aType AgentType, optFns ...func(o *Options)) (*Executor, error) {
	opts := Options{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	var (
		agent schema.Agent
		err   error
	)

	switch aType {
	case ReactDescriptionAgentType:
		agent, err = NewReactDescription(llm, tools)
		if err != nil {
			return nil, err
		}
	case ConversationalReactDescriptionAgentType:
		agent, err = NewConversationalReactDescription(llm, tools)
		if err != nil {
			return nil, err
		}
	case OpenAIFunctionsAgentType:
		agent, err = NewOpenAIFunctions(llm, tools)
		if err != nil {
			return nil, err
		}
	case ReactDocstoreAgentType:
		return nil, fmt.Errorf("agentType %s is not implemented", aType)
	default:
		return nil, ErrUnknownAgentType
	}

	return NewExecutor(agent, tools)
}

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
