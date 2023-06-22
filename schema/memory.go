package schema

import "context"

type Memory interface {
	// Input keys this memory class will load dynamically.
	MemoryKeys() []string
	// Return key-value pairs given the text input to the chain.
	// If None, return all memories
	LoadMemoryVariables(ctx context.Context, inputs map[string]any) (map[string]any, error)
	// Save the context of this model run to memory.
	SaveContext(ctx context.Context, inputs map[string]any, outputs map[string]any) error
	// Clear memory contents.
	Clear(ctx context.Context) error
}

type ChatMessageHistory interface {
	// Messages returns the messages stored in the store.
	Messages(ctx context.Context) (ChatMessages, error)
	// Add a user message to the store.
	AddUserMessage(ctx context.Context, text string) error
	// Add an AI message to the store.
	AddAIMessage(ctx context.Context, text string) error
	// Add a self-created message to the store.
	AddMessage(ctx context.Context, message ChatMessage) error
	// Remove all messages from the store.
	Clear(ctx context.Context) error
}
