package schema

import "context"

type LLMStartManagerInput struct {
	LLMType          string
	Prompt           string
	InvocationParams map[string]any
}

type LLMStartInput struct {
	*LLMStartManagerInput
	RunID string
}

type ChatModelStartManagerInput struct {
	ChatModelType    string
	Messages         ChatMessages
	InvocationParams map[string]any
}

type ChatModelStartInput struct {
	*ChatModelStartManagerInput
	RunID string
}

type ModelNewTokenManagerInput struct {
	Token string
}

type ModelNewTokenInput struct {
	*ModelNewTokenManagerInput
	RunID string
}

type ModelEndManagerInput struct {
	Result *ModelResult
}

type ModelEndInput struct {
	*ModelEndManagerInput
	RunID string
}

type ModelErrorManagerInput struct {
	Error error
}

type ModelErrorInput struct {
	*ModelErrorManagerInput
	RunID string
}

type ChainStartManagerInput struct {
	ChainType string
	Inputs    ChainValues
}

type ChainStartInput struct {
	*ChainStartManagerInput
	RunID string
}

type ChainEndManagerInput struct {
	Outputs ChainValues
}

type ChainEndInput struct {
	*ChainEndManagerInput
	RunID string
}

type ChainErrorManagerInput struct {
	Error error
}

type ChainErrorInput struct {
	*ChainErrorManagerInput
	RunID string
}

type AgentActionManagerInput struct {
	Action *AgentAction
}

type AgentActionInput struct {
	*AgentActionManagerInput
	RunID string
}

type AgentFinishManagerInput struct {
	Finish *AgentFinish
}

type AgentFinishInput struct {
	*AgentFinishManagerInput
	RunID string
}

type ToolStartManagerInput struct {
	ToolName string
	Input    *ToolInput
}

type ToolStartInput struct {
	*ToolStartManagerInput
	RunID string
}

type ToolEndManagerInput struct {
	Output string
}

type ToolEndInput struct {
	*ToolEndManagerInput
	RunID string
}

type ToolErrorManagerInput struct {
	Error error
}

type ToolErrorInput struct {
	*ToolErrorManagerInput
	RunID string
}

type TextManagerInput struct {
	Text string
}

type TextInput struct {
	*TextManagerInput
	RunID string
}

type RetrieverStartManagerInput struct {
	Query string
}

type RetrieverStartInput struct {
	*RetrieverStartManagerInput
	RunID string
}

type RetrieverEndManagerInput struct {
	Docs []Document
}

type RetrieverEndInput struct {
	*RetrieverEndManagerInput
	RunID string
}

type RetrieverErrorManagerInput struct {
	Error error
}

type RetrieverErrorInput struct {
	*RetrieverErrorManagerInput
	RunID string
}

type Callback interface {
	AlwaysVerbose() bool
	RaiseError() bool
	OnLLMStart(ctx context.Context, input *LLMStartInput) error
	OnChatModelStart(ctx context.Context, input *ChatModelStartInput) error
	OnModelNewToken(ctx context.Context, input *ModelNewTokenInput) error
	OnModelEnd(ctx context.Context, input *ModelEndInput) error
	OnModelError(ctx context.Context, input *ModelErrorInput) error
	OnChainStart(ctx context.Context, input *ChainStartInput) error
	OnChainEnd(ctx context.Context, input *ChainEndInput) error
	OnChainError(ctx context.Context, input *ChainErrorInput) error
	OnAgentAction(ctx context.Context, input *AgentActionInput) error
	OnAgentFinish(ctx context.Context, input *AgentFinishInput) error
	OnToolStart(ctx context.Context, input *ToolStartInput) error
	OnToolEnd(ctx context.Context, input *ToolEndInput) error
	OnToolError(ctx context.Context, input *ToolErrorInput) error
	OnText(ctx context.Context, input *TextInput) error
	OnRetrieverStart(ctx context.Context, input *RetrieverStartInput) error
	OnRetrieverEnd(ctx context.Context, input *RetrieverEndInput) error
	OnRetrieverError(ctx context.Context, input *RetrieverErrorInput) error
}

type CallbackManager interface {
	OnLLMStart(ctx context.Context, input *LLMStartManagerInput) (CallbackManagerForModelRun, error)
	OnChatModelStart(ctx context.Context, input *ChatModelStartManagerInput) (CallbackManagerForModelRun, error)
	OnChainStart(ctx context.Context, input *ChainStartManagerInput) (CallbackManagerForChainRun, error)
	OnToolStart(ctx context.Context, input *ToolStartManagerInput) (CallbackManagerForToolRun, error)
	OnRetrieverStart(ctx context.Context, input *RetrieverStartManagerInput) (CallbackManagerForRetrieverRun, error)
	RunID() string
}

type CallbackManagerForChainRun interface {
	OnChainEnd(ctx context.Context, input *ChainEndManagerInput) error
	OnChainError(ctx context.Context, input *ChainErrorManagerInput) error
	OnAgentAction(ctx context.Context, input *AgentActionManagerInput) error
	OnAgentFinish(ctx context.Context, input *AgentFinishManagerInput) error
	OnText(ctx context.Context, input *TextManagerInput) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallbackManagerForModelRun interface {
	OnModelNewToken(ctx context.Context, input *ModelNewTokenManagerInput) error
	OnModelEnd(ctx context.Context, input *ModelEndManagerInput) error
	OnModelError(ctx context.Context, input *ModelErrorManagerInput) error
	OnText(ctx context.Context, input *TextManagerInput) error
	GetInheritableCallbacks() []Callback
	RunID() string
}

type CallbackManagerForToolRun interface {
	OnToolEnd(ctx context.Context, input *ToolEndManagerInput) error
	OnToolError(ctx context.Context, input *ToolErrorManagerInput) error
	OnText(ctx context.Context, input *TextManagerInput) error
}

type CallbackManagerForRetrieverRun interface {
	OnRetrieverEnd(ctx context.Context, input *RetrieverEndManagerInput) error
	OnRetrieverError(ctx context.Context, input *RetrieverErrorManagerInput) error
}

type CallbackOptions struct {
	Callbacks []Callback
	Verbose   bool
}
