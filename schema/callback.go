package schema

type Callback interface {
	AlwaysVerbose() bool
	RaiseError() bool
	OnLLMStart(llmName string, prompts []string) error
	OnLLMNewToken(token string) error
	OnLLMEnd(result *LLMResult) error
	OnLLMError(llmError error) error
	OnChainStart(chainName string, inputs *ChainValues) error
	OnChainEnd(outputs *ChainValues) error
	OnChainError(chainError error) error
	OnToolStart(toolName string, input string) error
	OnToolEnd(output string) error
	OnToolError(toolError error) error
	// OnText() error
	// OnAgentAction() error
	// OnAgentFinish() error
}

type CallbackManager interface {
	OnLLMStart(llmName string, prompts []string) (CallBackManagerForLLMRun, error)
	OnChainStart(chainName string, inputs *ChainValues) (CallBackManagerForChainRun, error)
	OnToolStart(toolName string, input string) (CallBackManagerForToolRun, error)
	RunID() string
}

type CallBackManagerForChainRun interface {
	OnChainEnd(outputs *ChainValues) error
	OnChainError(chainError error) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallBackManagerForLLMRun interface {
	OnLLMNewToken(token string) error
	OnLLMEnd(result *LLMResult) error
	OnLLMError(llmError error) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallBackManagerForToolRun interface {
	OnToolEnd(output string) error
	OnToolError(toolError error) error
}

type CallbackOptions struct {
	Callbacks []Callback
	Verbose   bool
}
