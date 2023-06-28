package schema

type Callback interface {
	AlwaysVerbose() bool
	RaiseError() bool
	OnLLMStart(llmName string, prompts []string) error
	OnChatModelStart(chatModelName string, messages []ChatMessages) error
	OnModelNewToken(token string) error
	OnModelEnd(result ModelResult) error
	OnModelError(llmError error) error
	OnChainStart(chainName string, inputs ChainValues) error
	OnChainEnd(outputs ChainValues) error
	OnChainError(chainError error) error
	OnAgentAction(action AgentAction) error
	OnAgentFinish(finish AgentFinish) error
	OnToolStart(toolName string, input string) error
	OnToolEnd(output string) error
	OnToolError(toolError error) error
	OnText(text string) error
}

type CallbackManager interface {
	OnLLMStart(llmName string, prompts []string) (CallBackManagerForModelRun, error)
	OnChatModelStart(chatModelName string, messages []ChatMessages) (CallBackManagerForModelRun, error)
	OnChainStart(chainName string, inputs ChainValues) (CallBackManagerForChainRun, error)
	OnToolStart(toolName string, input string) (CallBackManagerForToolRun, error)
	RunID() string
}

type CallBackManagerForChainRun interface {
	OnChainEnd(outputs ChainValues) error
	OnChainError(chainError error) error
	OnAgentAction(action AgentAction) error
	OnAgentFinish(finish AgentFinish) error
	OnText(text string) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallBackManagerForModelRun interface {
	OnModelNewToken(token string) error
	OnModelEnd(result ModelResult) error
	OnModelError(llmError error) error
	OnText(text string) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallBackManagerForToolRun interface {
	OnToolEnd(output string) error
	OnToolError(toolError error) error
	OnText(text string) error
}

type CallbackOptions struct {
	Callbacks []Callback
	Verbose   bool
}
