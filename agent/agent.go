package agent

import "github.com/hupe1980/golc"

type AgentType string

const (
	ZeroShotReactDescriptionAgentType AgentType = "zero-shot-react-description"
)

type ZeroShotReactDescriptionAgent struct{}

func NewZeroShotReactDescriptionAgent(llm golc.LLM, tools []golc.Tool) *ZeroShotReactDescriptionAgent {
	return &ZeroShotReactDescriptionAgent{}
}
