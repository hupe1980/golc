package schema

type Memory interface {
	// Input keys this memory class will load dynamically.
	MemoryVariables() []string
	// Return key-value pairs given the text input to the chain.
	// If None, return all memories
	LoadMemoryVariables(inputs map[string]any) (map[string]any, error)
	// Save the context of this model run to memory.
	SaveContext(inputs map[string]any, outputs map[string]any) error
	// Clear memory contents.
	Clear() error
}

type ChatMessageHistory interface {
	// Messages returns the messages stored in the store.
	Messages() ([]ChatMessage, error)
	// Add a user message to the store.
	AddUserMessage(text string) error
	// Add an AI message to the store.
	AddAIMessage(text string) error
	// Add a self-created message to the store.
	AddMessage(message ChatMessage) error
	// Remove all messages from the store.
	Clear() error
}
