package agent

import "github.com/hupe1980/golc"

func Initialize(llm golc.LLM, tools []golc.Tool, aType AgentType) (*Executor, error) {
	var agent golc.Agent

	switch aType {
	case ZeroShotReactDescriptionAgentType:
		agent = NewZeroShotReactDescriptionAgent(llm, tools)
	default:
		return nil, ErrUnknownAgentType
	}

	return NewExecutor(agent)
}
