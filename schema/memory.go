package schema

import "context"

// Memory is an interface for managing context and memory in a chain.
type Memory interface {
	// MemoryKeys returns the memory keys.
	MemoryKeys() []string
	// LoadMemoryVariables returns key-value pairs given the text input to the chain.
	LoadMemoryVariables(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
	// SaveContext saves the context of this model run to memory.
	SaveContext(ctx context.Context, inputs map[string]interface{}, outputs map[string]interface{}) error
	// Clear clears memory contents.
	Clear(ctx context.Context) error
}

type ChatMessageHistory interface {
	// Messages returns the messages stored in the store.
	Messages(ctx context.Context) (ChatMessages, error)
	// AddUserMessage adds a user message to the store.
	AddUserMessage(ctx context.Context, text string) error
	// AddAIMessage adds an AI message to the store.
	AddAIMessage(ctx context.Context, text string) error
	// AddMessage adds a self-created message to the store.
	AddMessage(ctx context.Context, message ChatMessage) error
	// Clear removes all messages from the store.
	Clear(ctx context.Context) error
}
