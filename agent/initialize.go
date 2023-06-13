package agent

import "github.com/hupe1980/golc"

func Initialize(llm golc.LLM, tools []golc.Tool, aType AgentType) (*Executor, error) {
	var (
		agent golc.Agent
		err   error
	)

	switch aType {
	case ZeroShotReactDescriptionAgentType:
		agent, err = NewZeroShotReactDescriptionAgent(llm, tools)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnknownAgentType
	}

	return NewExecutor(agent)
}
