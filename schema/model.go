package schema

import (
	"context"

	"github.com/hupe1980/golc/integration/jsonschema"
)

// Generation represents a generated text along with its corresponding chat message and additional information.
type Generation struct {
	Text    string
	Message ChatMessage
	Info    map[string]any
}

// ModelResult represents the result of a model generation.
type ModelResult struct {
	Generations []Generation
	LLMOutput   map[string]any
}

// PromptValue is an interface representing a prompt value for LLMs and chat models.
type PromptValue interface {
	// String returns the string representation of the prompt value.
	String() string

	// Messages returns the chat messages associated with the prompt value.
	Messages() ChatMessages
}

// PromptTemplate is an interface for creating templates that can be formatted with dynamic values.
type PromptTemplate interface {
	// Format applies values to the template and returns the formatted result as a string.
	Format(values map[string]any) (string, error)

	// FormatPrompt applies values to the template and returns a PromptValue representation of the formatted result.
	FormatPrompt(values map[string]any) (PromptValue, error)

	// InputVariables returns a list of input variables used in the template.
	InputVariables() []string

	// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
	OutputParser() (OutputParser[any], bool)
}

// Tokenizer is an interface for tokenizing text.
type Tokenizer interface {
	// GetNumTokens returns the number of tokens in the provided text.
	GetNumTokens(ctx context.Context, text string) (uint, error)
	// GetNumTokensFromMessage returns the number of tokens in the provided chat messages.
	GetNumTokensFromMessage(ctx context.Context, messages ChatMessages) (uint, error)
}

type FunctionDefinitionParameters struct {
	Type       string                        `json:"type"`
	Properties map[string]*jsonschema.Schema `json:"properties"`
	Required   []string                      `json:"required"`
}

type FunctionDefinition struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
	Parameters  FunctionDefinitionParameters `json:"parameters"`
}

type GenerateOptions struct {
	CallbackManger    CallbackManagerForModelRun
	Stop              []string
	Functions         []FunctionDefinition
	ForceFunctionCall bool
}

// LLM is the interface for language models.
type LLM interface {
	Model
	// Generate generates text based on the provided prompt and options.
	Generate(ctx context.Context, prompt string, optFns ...func(o *GenerateOptions)) (*ModelResult, error)
}

// ChatModel is the interface for chat models.
type ChatModel interface {
	Model
	// Generate generates text based on the provided chat messages and options.
	Generate(ctx context.Context, messages ChatMessages, optFns ...func(o *GenerateOptions)) (*ModelResult, error)
}

// Model is the interface for language models and chat models.
type Model interface {
	Tokenizer
	// Type returns the type of the model.
	Type() string
	// Verbose returns the verbosity setting of the model.
	Verbose() bool
	// Callbacks returns the registered callbacks of the model.
	Callbacks() []Callback
	// InvocationParams returns the parameters used in the model invocation.
	InvocationParams() map[string]any
}

// Embedder is the interface for creating vector embeddings from texts.
type Embedder interface {
	// BatchEmbedText embeds a list of texts and returns their embeddings.
	BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error)
	// EmbedText embeds a single text and returns its embedding.
	EmbedText(ctx context.Context, text string) ([]float32, error)
}

// OutputParser is an interface for parsing the output of an LLM call.
type OutputParser[T any] interface {
	// Parse parses the output of an LLM call.
	ParseResult(result Generation) (any, error)
	// Parse parses the output of an LLM call.
	Parse(text string) (T, error)
	// ParseWithPrompt parses the output of an LLM call with the prompt used.
	ParseWithPrompt(text string, prompt PromptValue) (T, error)
	// GetFormatInstructions returns a string describing the format of the output.
	GetFormatInstructions() string
	// Type returns the string type key uniquely identifying this class of parser
	Type() string
}
