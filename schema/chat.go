package schema

import (
	"fmt"
	"strings"
)

// FunctionCall represents a function call with optional arguments in JSON format.
type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

// ChatMessageType represents the type of a chat message.
type ChatMessageType string

const (
	ChatMessageTypeHuman    ChatMessageType = "human"
	ChatMessageTypeAI       ChatMessageType = "ai"
	ChatMessageTypeSystem   ChatMessageType = "system"
	ChatMessageTypeGeneric  ChatMessageType = "generic"
	ChatMessageTypeFunction ChatMessageType = "function"
)

// ChatMessageExtension represents additional data associated with a chat message.
type ChatMessageExtension struct {
	FunctionCall *FunctionCall `json:"functionCall,omitempty"`
}

// ChatMessage is an interface for different types of chat messages.
type ChatMessage interface {
	// Content returns the content of the chat message.
	Content() string
	// Type returns the type of the chat message.
	Type() ChatMessageType
}

// ChatMessageToMap converts a ChatMessage to a map representation.
func ChatMessageToMap(cm ChatMessage) map[string]string {
	m := map[string]string{
		"type":    string(cm.Type()),
		"content": cm.Content(),
	}

	if fm, ok := cm.(*FunctionChatMessage); ok {
		m["name"] = fm.Name()
	} else if gm, ok := cm.(*GenericChatMessage); ok {
		m["role"] = gm.Role()
	}

	return m
}

// MapToChatMessage converts a map representation back to a ChatMessage.
func MapToChatMessage(m map[string]string) (ChatMessage, error) {
	switch ChatMessageType(m["type"]) {
	case ChatMessageTypeHuman:
		return NewHumanChatMessage(m["content"]), nil
	case ChatMessageTypeAI:
		return NewAIChatMessage(m["content"]), nil
	case ChatMessageTypeSystem:
		return NewSystemChatMessage(m["content"]), nil
	case ChatMessageTypeGeneric:
		return NewGenericChatMessage(m["content"], m["role"]), nil
	case ChatMessageTypeFunction:
		return NewFunctionChatMessage(m["content"], m["name"]), nil
	default:
		return nil, fmt.Errorf("unknown chat message type: %s", m["type"])
	}
}

// HumanChatMessage represents a chat message from a human.
type HumanChatMessage struct {
	content string
}

// NewHumanChatMessage creates a new HumanChatMessage instance.
func NewHumanChatMessage(content string) *HumanChatMessage {
	return &HumanChatMessage{
		content: content,
	}
}

// Type returns the type of the chat message.
func (m HumanChatMessage) Type() ChatMessageType { return ChatMessageTypeHuman }

// Content returns the content of the chat message.
func (m HumanChatMessage) Content() string { return m.content }

// AIChatMessage represents a chat message from an AI.
type AIChatMessage struct {
	content string
	ext     ChatMessageExtension
}

// NewAIChatMessage creates a new AIChatMessage instance.
func NewAIChatMessage(content string, extFns ...func(o *ChatMessageExtension)) *AIChatMessage {
	ext := ChatMessageExtension{}

	for _, fn := range extFns {
		fn(&ext)
	}

	return &AIChatMessage{
		content: content,
		ext:     ext,
	}
}

// Type returns the type of the chat message.
func (m AIChatMessage) Type() ChatMessageType { return ChatMessageTypeAI }

// Content returns the content of the chat message.
func (m AIChatMessage) Content() string { return m.content }

// Extension returns the extension data of the chat message.
func (m AIChatMessage) Extension() ChatMessageExtension { return m.ext }

// SystemChatMessage represents a chat message from the system.
type SystemChatMessage struct {
	content string
}

// NewSystemChatMessage creates a new SystemChatMessage instance.
func NewSystemChatMessage(content string) *SystemChatMessage {
	return &SystemChatMessage{
		content: content,
	}
}

// Type returns the type of the chat message.
func (m SystemChatMessage) Type() ChatMessageType { return ChatMessageTypeSystem }

// Content returns the content of the chat message.
func (m SystemChatMessage) Content() string { return m.content }

// GenericChatMessage represents a generic chat message with an associated role.
type GenericChatMessage struct {
	content string
	role    string
}

// NewGenericChatMessage creates a new GenericChatMessage instance.
func NewGenericChatMessage(content, role string) *GenericChatMessage {
	return &GenericChatMessage{
		content: content,
		role:    role,
	}
}

// Type returns the type of the chat message.
func (m GenericChatMessage) Type() ChatMessageType { return ChatMessageTypeGeneric }

// Content returns the content of the chat message.
func (m GenericChatMessage) Content() string { return m.content }

// Role returns the role associated with the chat message.
func (m GenericChatMessage) Role() string { return m.role }

// FunctionChatMessage represents a chat message for a function call.
type FunctionChatMessage struct {
	name    string
	content string
}

// NewFunctionChatMessage creates a new FunctionChatMessage instance.
func NewFunctionChatMessage(name, content string) *FunctionChatMessage {
	return &FunctionChatMessage{
		name:    name,
		content: content,
	}
}

// Type returns the type of the chat message.
func (m FunctionChatMessage) Type() ChatMessageType { return ChatMessageTypeFunction }

// Content returns the content of the chat message.
func (m FunctionChatMessage) Content() string { return m.content }

// Name returns the name of the function associated with the chat message.
func (m FunctionChatMessage) Name() string { return m.name }

// ChatMessages represents a slice of ChatMessage.
type ChatMessages []ChatMessage

// StringifyChatMessagesOptions represents options for formatting ChatMessages.
type StringifyChatMessagesOptions struct {
	HumanPrefix    string
	AIPrefix       string
	SystemPrefix   string
	FunctionPrefix string
}

// Format formats the ChatMessages into a single string representation.
func (cm ChatMessages) Format(optFns ...func(o *StringifyChatMessagesOptions)) (string, error) {
	opts := StringifyChatMessagesOptions{
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		SystemPrefix:   "System",
		FunctionPrefix: "Function",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	result := []string{}

	for _, message := range cm {
		var role string

		switch message.Type() {
		case ChatMessageTypeHuman:
			role = opts.HumanPrefix
		case ChatMessageTypeAI:
			role = opts.AIPrefix
		case ChatMessageTypeSystem:
			role = opts.SystemPrefix
		case ChatMessageTypeGeneric:
			role = message.(*GenericChatMessage).Role()
		case ChatMessageTypeFunction:
			role = opts.FunctionPrefix
		default:
			return "", fmt.Errorf("unknown chat message type: %s", message.Type())
		}

		result = append(result, fmt.Sprintf("%s: %s", role, message.Content()))
	}

	return strings.Join(result, "\n"), nil
}
