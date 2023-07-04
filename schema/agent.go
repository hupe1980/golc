package schema

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
)

type ToolInput struct {
	sinput     string
	structured bool
}

func NewToolInputFromString(input string) *ToolInput {
	return &ToolInput{
		sinput:     input,
		structured: false,
	}
}

func NewToolInputFromArguments(input string) *ToolInput {
	return &ToolInput{
		sinput:     input,
		structured: true,
	}
}

func (ti *ToolInput) Structured() bool {
	return ti.structured
}

func (ti *ToolInput) GetString() (string, error) {
	if ti.structured {
		return "", errors.New("cannot return string for strutured input")
	}

	return ti.sinput, nil
}

func (ti *ToolInput) Unmarshal(args any) error {
	if ptr, ok := args.(*string); ok {
		temp := make(map[string]string, 1)
		if err := json.Unmarshal([]byte(ti.sinput), &temp); err != nil {
			return err
		}

		*ptr = temp["__arg1"]

		return nil
	}

	if err := json.Unmarshal([]byte(ti.sinput), args); err != nil {
		return err
	}

	return nil
}

func (ti *ToolInput) String() string {
	return ti.sinput
}

// AgentAction is the agent's action to take.
type AgentAction struct {
	Tool       string
	ToolInput  *ToolInput
	Log        string
	MessageLog ChatMessages
}

// AgentStep is a step of the agent.
type AgentStep struct {
	Action      *AgentAction
	Observation string
}

// AgentFinish is the agent's return value.
type AgentFinish struct {
	ReturnValues map[string]any
	Log          string
}

type Agent interface {
	Plan(ctx context.Context, intermediateSteps []AgentStep, inputs map[string]string) ([]*AgentAction, *AgentFinish, error)
	InputKeys() []string
	OutputKeys() []string
}

type Tool interface {
	Name() string
	Description() string
	Run(ctx context.Context, input any) (string, error)
	ArgsType() reflect.Type
}
