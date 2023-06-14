package agent

import "errors"

var (
	ErrUnknownAgentType       = errors.New("unknown agent type")
	ErrAgentNoReturn          = errors.New("no actions or finish was returned by the agent")
	ErrNotFinished            = errors.New("agent not finished before max iterations")
	ErrExecutorInputNotString = errors.New("input to executor is not a string")
	ErrInvalidChainReturnType = errors.New("agent chain did not return a string")
	ErrUnableToParseOutput    = errors.New("unable to parse agent output")
)
