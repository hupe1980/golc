package schema

import "context"

type Generation struct {
	Text    string
	Message ChatMessage
	Info    map[string]any
}

type LLMResult struct {
	Generations [][]Generation
	LLMOutput   map[string]any
}

type ChainValues map[string]any

type CallOptions struct {
	CallbackManger CallBackManagerForChainRun
	Stop           []string
}

type Chain interface {
	// Call executes the chain with the given context and inputs.
	// It returns the outputs of the chain or an error, if any.
	Call(ctx context.Context, inputs ChainValues, optFns ...func(o *CallOptions)) (ChainValues, error)
	// Type returns the type of the chain.
	Type() string
	// Verbose returns the verbosity setting of the chain.
	Verbose() bool
	// Callbacks returns the callbacks associated with the chain.
	Callbacks() []Callback
	// Memory returns the memory associated with the chain.
	Memory() Memory
	// InputKeys returns the expected input keys.
	InputKeys() []string
	// OutputKeys returns the output keys the chain will return.
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

type GenerateOptions struct {
	CallbackManger CallBackManagerForLLMRun
	Stop           []string
}

type LLM interface {
	Model
	Generate(ctx context.Context, prompts []string, optFns ...func(o *GenerateOptions)) (*LLMResult, error)
}

type ChatModel interface {
	Model
	Generate(ctx context.Context, messages ChatMessages) (*LLMResult, error)
}

type Model interface {
	Tokenizer
	Type() string
	Verbose() bool
	Callbacks() []Callback
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
	ParseResult(result []Generation) (any, error)
	// Parse parses the output of an LLM call.
	Parse(text string) (T, error)
	// ParseWithPrompt parses the output of an LLM call with the prompt used.
	ParseWithPrompt(text string, prompt PromptValue) (T, error)
	// GetFormatInstructions returns a string describing the format of the output.
	GetFormatInstructions() (string, error)
	// Type returns the string type key uniquely identifying this class of parser
	Type() string
}
