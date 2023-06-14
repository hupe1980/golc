package agent

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
)

// Compile time check to ensure Executor satisfies the chain interface.
var _ golc.Chain = (*Executor)(nil)

type ExecutorOptions struct {
	MaxIterations int
}

type Executor struct {
	agent    golc.Agent
	toolsMap map[string]golc.Tool
	opts     ExecutorOptions
}

func NewExecutor(agent golc.Agent, tools []golc.Tool) (*Executor, error) {
	opts := ExecutorOptions{
		MaxIterations: 5,
	}

	// Construct a mapping of tool name to tool for easy lookup
	toolsMap := make(map[string]golc.Tool, len(tools))
	for _, tool := range tools {
		toolsMap[tool.Name()] = tool
	}

	return &Executor{
		agent:    agent,
		toolsMap: toolsMap,
		opts:     opts,
	}, nil
}

func (e Executor) Call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	inputs, err := inputsToString(values)
	if err != nil {
		return nil, err
	}

	steps := []golc.AgentStep{}

	for i := 0; i <= e.opts.MaxIterations; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			actions, finish, err := e.agent.Plan(ctx, steps, inputs)
			if err != nil {
				return nil, err
			}

			if len(actions) == 0 && finish == nil {
				return nil, ErrAgentNoReturn
			}

			if finish != nil {
				return finish.ReturnValues, nil
			}

			for _, action := range actions {
				tool, ok := e.toolsMap[action.Tool]
				if !ok {
					steps = append(steps, golc.AgentStep{
						Action:      action,
						Observation: fmt.Sprintf("%s is not a valid tool, try another one", action.Tool),
					})

					continue
				}

				observation, err := tool.Run(ctx, action.ToolInput)
				if err != nil {
					return nil, err
				}

				steps = append(steps, golc.AgentStep{
					Action:      action,
					Observation: observation,
				})
			}
		}
	}

	return nil, ErrNotFinished
}

func (e Executor) InputKeys() []string {
	return e.agent.InputKeys()
}

func (e Executor) OutputKeys() []string {
	return e.agent.OutputKeys()
}
