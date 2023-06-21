package zep

// Message represents a message in a conversation.
type Message struct {
	// The content of the message.
	Content string `json:"content,omitempty"`
	// The timestamp of when the message was created.
	CreatedAt string `json:"created_at,omitempty"`
	// Metadata associated with the message.
	Metadata interface{} `json:"metadata,omitempty"`
	// The role of the sender of the message (e.g., "user", "assistant").
	Role string `json:"role,omitempty"`
	// The number of tokens in the message.
	TokenCount int `json:"token_count,omitempty"`
	// The unique identifier of the message.
	UUID string `json:"uuid,omitempty"`
}

// SearchPayload represents a search payload for querying memory.
type SearchPayload struct {
	// Metadata associated with the search query.
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The text of the search query.
	Text string `json:"text,omitempty"`
}

// SearchResult represents a search result from querying memory.
type SearchResult struct {
	// The distance metric of the search result.
	Dist float64 `json:"dist,omitempty"`
	// The message associated with the search result.
	Message Message `json:"message,omitempty"`
	// Metadata associated with the search result.
	Meta interface{} `json:"meta,omitempty"`
	// The summary of the search result.
	Summary Summary `json:"summary,omitempty"`
}

// Summary represents a summary of a conversation.
type Summary struct {
	// The content of the summary.
	Content string `json:"content,omitempty"`
	// The timestamp of when the summary was created.
	CreatedAt string `json:"created_at,omitempty"`
	// // Metadata associated with the summary.
	Metadata interface{} `json:"metadata,omitempty"`
	// The unique identifier of the most recent message in the conversation.
	RecentMessageUUID string `json:"recent_message_uuid,omitempty"`
	// The number of tokens in the summary.
	TokenCount int `json:"token_count,omitempty"`
	// The unique identifier of the summary.
	UUID string `json:"uuid,omitempty"`
}

// Represents a memory object with messages, metadata, and other attributes.
type Memory struct {
	// A list of message objects, where each message contains a role and content.
	Messages []Message `json:"messages,omitempty"`
	// A dictionary containing metadata associated with the memory.
	Metadata interface{} `json:"metadata,omitempty"`
	// A Summary object.
	Summary Summary `json:"summary,omitempty"`
	// A unique identifier for the memory.
	UUID string `json:"uuid,omitempty"`
	// The timestamp when the memory was created.
	CreatedAt string `json:"created_at,omitempty"`
	// The token count of the memory.
	TokenCount int `json:"token_count,omitempty"`
}

type APIError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
