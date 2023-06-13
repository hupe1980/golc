package agent

import (
	"context"

	"github.com/hupe1980/golc"
)

type Executor struct {
	agent golc.Agent
}

func NewExecutor(agent golc.Agent) (*Executor, error) {
	return &Executor{
		agent: agent,
	}, nil
}

func (e Executor) Plan(ctx context.Context) {
}

func (e Executor) InputKeys() []string {
	return e.agent.InputKeys()
}

func (e Executor) OutputKeys() []string {
	return e.agent.OutputKeys()
}
