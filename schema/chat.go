package schema

import (
	"fmt"
	"strings"
)

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type ChatMessageType string

const (
	ChatMessageTypeHuman    ChatMessageType = "human"
	ChatMessageTypeAI       ChatMessageType = "ai"
	ChatMessageTypeSystem   ChatMessageType = "system"
	ChatMessageTypeGeneric  ChatMessageType = "generic"
	ChatMessageTypeFunction ChatMessageType = "function"
)

type ChatMessageExtension struct {
	FunctionCall *FunctionCall `json:"functionCall,omitempty"`
}

type ChatMessage interface {
	Content() string
	Type() ChatMessageType
}

func ChatMessageToMap(cm ChatMessage) map[string]string {
	m := map[string]string{
		"type":    string(cm.Type()),
		"content": cm.Content(),
	}

	if gm, ok := cm.(GenericChatMessage); ok {
		m["role"] = gm.Role()
	}

	return m
}

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

type HumanChatMessage struct {
	content string
}

func NewHumanChatMessage(content string) *HumanChatMessage {
	return &HumanChatMessage{
		content: content,
	}
}

func (m HumanChatMessage) Type() ChatMessageType { return ChatMessageTypeHuman }
func (m HumanChatMessage) Content() string       { return m.content }

type AIChatMessage struct {
	content string
	ext     ChatMessageExtension
}

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

func (m AIChatMessage) Type() ChatMessageType           { return ChatMessageTypeAI }
func (m AIChatMessage) Content() string                 { return m.content }
func (m AIChatMessage) Extension() ChatMessageExtension { return m.ext }

type SystemChatMessage struct {
	content string
}

func NewSystemChatMessage(content string) *SystemChatMessage {
	return &SystemChatMessage{
		content: content,
	}
}

func (m SystemChatMessage) Type() ChatMessageType { return ChatMessageTypeSystem }
func (m SystemChatMessage) Content() string       { return m.content }

type GenericChatMessage struct {
	content string
	role    string
}

func NewGenericChatMessage(content, role string) *GenericChatMessage {
	return &GenericChatMessage{
		content: content,
		role:    role,
	}
}

func (m GenericChatMessage) Type() ChatMessageType { return ChatMessageTypeGeneric }
func (m GenericChatMessage) Content() string       { return m.content }
func (m GenericChatMessage) Role() string          { return m.role }

type FunctionChatMessage struct {
	name    string
	content string
}

func NewFunctionChatMessage(name, content string) *FunctionChatMessage {
	return &FunctionChatMessage{
		name:    name,
		content: content,
	}
}

func (m FunctionChatMessage) Type() ChatMessageType { return ChatMessageTypeFunction }
func (m FunctionChatMessage) Content() string       { return m.content }
func (m FunctionChatMessage) Name() string          { return m.name }

type ChatMessages []ChatMessage

type StringifyChatMessagesOptions struct {
	HumanPrefix    string
	AIPrefix       string
	SystemPrefix   string
	FunctionPrefix string
}

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
			role = message.(GenericChatMessage).Role()
		case ChatMessageTypeFunction:
			role = opts.FunctionPrefix
		default:
			return "", fmt.Errorf("unknown chat message type: %s", message.Type())
		}

		result = append(result, fmt.Sprintf("%s: %s", role, message.Content()))
	}

	return strings.Join(result, "\n"), nil
}
