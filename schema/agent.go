package schema

import "context"

// AgentAction is the agent's action to take.
type AgentAction struct {
	Tool      string
	ToolInput string
	Log       string
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
	Run(ctx context.Context, query string) (string, error)
}
