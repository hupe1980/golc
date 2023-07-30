package agent

import "errors"

var (
	ErrAgentNoReturn          = errors.New("no actions or finish was returned by the agent")
	ErrNotFinished            = errors.New("agent not finished before max iterations")
	ErrInvalidChainReturnType = errors.New("agent chain did not return a string")
	ErrUnableToParseOutput    = errors.New("unable to parse agent output")
)
