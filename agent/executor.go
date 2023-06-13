package agent

import "github.com/hupe1980/golc"

type Executor struct {
	agent golc.Agent
}

func NewExecutor(agent golc.Agent) (*Executor, error) {
	return &Executor{
		agent: agent,
	}, nil
}
