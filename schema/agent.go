package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ToolInput represents an input for a tool, which can be either a structured input or a plain string input.
type ToolInput struct {
	sinput     string
	structured bool
}

// NewToolInputFromString creates a new ToolInput from a plain string input.
func NewToolInputFromString(input string) *ToolInput {
	return &ToolInput{
		sinput:     input,
		structured: false,
	}
}

// NewToolInputFromArguments creates a new ToolInput from a structured input.
func NewToolInputFromArguments(input string) *ToolInput {
	return &ToolInput{
		sinput:     input,
		structured: true,
	}
}

// Structured returns true if the ToolInput is a structured input; otherwise, it returns false.
func (ti *ToolInput) Structured() bool {
	return ti.structured
}

// GetString returns the plain string value from the ToolInput if it is not a structured input.
// If the ToolInput is structured, it returns an error.
func (ti *ToolInput) GetString() (string, error) {
	if ti.structured {
		return "", errors.New("cannot return string for structured input")
	}

	return ti.sinput, nil
}

// Unmarshal unmarshals the ToolInput into the provided argument, which should be a pointer to a valid data type.
// If the ToolInput is a structured input, it attempts to unmarshal it into a map with a single key "__arg1".
func (ti *ToolInput) Unmarshal(args any) error {
	if !ti.structured {
		if ptr, ok := args.(*string); ok {
			*ptr = ti.sinput
			return nil
		} else {
			return fmt.Errorf("cannot unmarshal value: %s", ti.sinput)
		}
	}

	if ptr, ok := args.(*string); ok {
		temp := make(map[string]string, 1)
		if err := json.Unmarshal([]byte(ti.sinput), &temp); err != nil {
			return err
		}

		*ptr = temp["__arg1"]

		return nil
	}

	return json.Unmarshal([]byte(ti.sinput), args)
}

// String returns the string representation of the ToolInput.
func (ti *ToolInput) String() string {
	return ti.sinput
}

// AgentAction represents an action that the agent will take.
type AgentAction struct {
	// Name of the tool to use for the action.
	Tool string
	// Input for the tool action.
	ToolInput *ToolInput
	// Log message associated with the action.
	Log string
	// Message log associated with the action.
	MessageLog ChatMessages
}

// AgentStep represents a step in the agent's action plan.
type AgentStep struct {
	// Action to be taken by the agent.
	Action *AgentAction
	// Observation made during the step.
	Observation string
}

// AgentFinish represents the return value of the agent.
type AgentFinish struct {
	// Return values from the agent.
	ReturnValues map[string]any
	Log          string
}

// Agent is an interface that defines the behavior of an agent.
type Agent interface {
	// Plan plans the agent's action given the intermediate steps and inputs.
	Plan(ctx context.Context, intermediateSteps []AgentStep, inputs ChainValues) ([]*AgentAction, *AgentFinish, error)
	// InputKeys returns the keys for expected input values for the agent.
	InputKeys() []string
	// OutputKeys returns the keys for the agent's output values.
	OutputKeys() []string
}

// Tool is an interface that defines the behavior of a tool.
type Tool interface {
	// Name returns the name of the tool.
	Name() string
	// Description returns the description of the tool.
	Description() string
	// Run executes the tool with the given input and returns the output.
	Run(ctx context.Context, input any) (string, error)
	// ArgsType returns the type of the input argument expected by the tool.
	ArgsType() reflect.Type
	// Verbose returns the verbosity setting of the tool.
	Verbose() bool
	// Callbacks returns the registered callbacks of the tool.
	Callbacks() []Callback
}
