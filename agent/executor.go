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

// ExecutorOptions holds configuration options for the Executor.
type ExecutorOptions struct {
	*schema.CallbackOptions
	MaxIterations int
	Memory        schema.Memory
}

// Executor represents an agent executor that executes a chain of actions based on inputs and a defined agent model.
type Executor struct {
	agent    schema.Agent
	toolsMap map[string]schema.Tool
	opts     ExecutorOptions
}

// NewExecutor creates a new instance of the Executor with the given agent and a list of available tools.
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

// Call executes the AgentExecutor chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (e Executor) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	// strInputs, err := inputsToString(inputs)
	// if err != nil {
	// 	return nil, err
	// }

	steps := []schema.AgentStep{}

	for i := 0; i <= e.opts.MaxIterations; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			actions, finish, err := e.agent.Plan(ctx, steps, inputs.Clone())
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

// Memory returns the memory associated with the chain.
func (e Executor) Memory() schema.Memory {
	return e.opts.Memory
}

// Type returns the type of the chain.
func (e Executor) Type() string {
	return "AgentExecutor"
}

// Verbose returns the verbosity setting of the chain.
func (e Executor) Verbose() bool {
	return e.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (e Executor) Callbacks() []schema.Callback {
	return e.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (e Executor) InputKeys() []string {
	return e.agent.InputKeys()
}

// OutputKeys returns the output keys the chain will return.
func (e Executor) OutputKeys() []string {
	return e.agent.OutputKeys()
}
