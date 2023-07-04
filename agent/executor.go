package agent

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
)

// Compile time check to ensure Executor satisfies the chain interface.
var _ schema.Chain = (*Executor)(nil)

type ExecutorOptions struct {
	*schema.CallbackOptions
	MaxIterations int
	Memory        schema.Memory
}

type Executor struct {
	agent    schema.Agent
	toolsMap map[string]schema.Tool
	opts     ExecutorOptions
}

func NewExecutor(agent schema.Agent, tools []schema.Tool) (*Executor, error) {
	opts := ExecutorOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		MaxIterations: 5,
	}

	// Construct a mapping of tool name to tool for easy lookup
	toolsMap := make(map[string]schema.Tool, len(tools))
	for _, tool := range tools {
		toolsMap[tool.Name()] = tool
	}

	return &Executor{
		agent:    agent,
		toolsMap: toolsMap,
		opts:     opts,
	}, nil
}

func (e Executor) Call(ctx context.Context, values schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	inputs, err := inputsToString(values)
	if err != nil {
		return nil, err
	}

	steps := []schema.AgentStep{}

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
				if cbErr := opts.CallbackManger.OnAgentFinish(ctx, &schema.AgentFinishManagerInput{
					Finish: finish,
				}); cbErr != nil {
					return nil, cbErr
				}

				return finish.ReturnValues, nil
			}

			for _, action := range actions {
				if cbErr := opts.CallbackManger.OnAgentAction(ctx, &schema.AgentActionManagerInput{
					Action: action,
				}); cbErr != nil {
					return nil, cbErr
				}

				t, ok := e.toolsMap[action.Tool]
				if !ok {
					steps = append(steps, schema.AgentStep{
						Action:      action,
						Observation: fmt.Sprintf("%s is not a valid tool, try another one", action.Tool),
					})

					continue
				}

				observation, err := tool.Run(ctx, t, action.ToolInput)
				if err != nil {
					return nil, err
				}

				steps = append(steps, schema.AgentStep{
					Action:      action,
					Observation: observation,
				})
			}
		}
	}

	return nil, ErrNotFinished
}

func (e Executor) Memory() schema.Memory {
	return e.opts.Memory
}

func (e Executor) Type() string {
	return "AgentExecutor"
}

func (e Executor) Verbose() bool {
	return e.opts.Verbose
}

func (e Executor) Callbacks() []schema.Callback {
	return e.opts.Callbacks
}

func (e Executor) InputKeys() []string {
	return e.agent.InputKeys()
}

func (e Executor) OutputKeys() []string {
	return e.agent.OutputKeys()
}
