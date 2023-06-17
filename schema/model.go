package schema

import "context"

type Generation struct {
	Text    string
	Message ChatMessage
	Info    map[string]any
}

type LLMResult struct {
	Generations [][]*Generation
	LLMOutput   map[string]any
}

type ChainValues map[string]any

type Chain interface {
	Call(ctx context.Context, inputs ChainValues) (ChainValues, error)
	InputKeys() []string
	OutputKeys() []string
}

type PromptValue interface {
	String() string
	Messages() ChatMessages
}

type Tokenizer interface {
	GetTokenIDs(text string) ([]int, error)
	GetNumTokens(text string) (int, error)
	GetNumTokensFromMessage(messages ChatMessages) (int, error)
}

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
	// OnToolStart() error
	// OnToolEnd() error
	// OnToolError() error
	// OnText() error
	// OnAgentAction() error
	// OnAgentFinish() error
}

type GenerateOptions struct {
	Stop      []string
	Callbacks []Callback
}

type LLM interface {
	Tokenizer
	GeneratePrompt(ctx context.Context, promptValues []PromptValue, optFns ...func(o *GenerateOptions)) (*LLMResult, error)
	Predict(ctx context.Context, text string, optFns ...func(o *GenerateOptions)) (string, error)
	PredictMessages(ctx context.Context, messages ChatMessages, optFns ...func(o *GenerateOptions)) (ChatMessage, error)
}

// Embedder is the interface for creating vector embeddings from texts.
type Embedder interface {
	// EmbedDocuments returns a vector for each text.
	EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error)
	// EmbedQuery embeds a single text.
	EmbedQuery(ctx context.Context, text string) ([]float64, error)
}

// OutputParser is an interface for parsing the output of an LLM call.
type OutputParser[T any] interface {
	// Parse parses the output of an LLM call.
	Parse(text string) (T, error)
	// ParseWithPrompt parses the output of an LLM call with the prompt used.
	ParseWithPrompt(text string, prompt PromptValue) (T, error)
	// GetFormatInstructions returns a string describing the format of the output.
	GetFormatInstructions() (string, error)
	// Type returns the string type key uniquely identifying this class of parser
	Type() string
}
