package schema

import (
	"fmt"
	"strings"
)

type ChatMessageType string

const (
	ChatMessageTypeHuman   ChatMessageType = "human"
	ChatMessageTypeAI      ChatMessageType = "ai"
	ChatMessageTypeSystem  ChatMessageType = "system"
	ChatMessageTypeGeneric ChatMessageType = "generic"
)

type ChatMessage interface {
	Text() string
	Type() ChatMessageType
}

func ChatMessageToMap(cm ChatMessage) map[string]string {
	m := map[string]string{
		"type": string(cm.Type()),
		"text": cm.Text(),
	}

	if gm, ok := cm.(GenericChatMessage); ok {
		m["role"] = gm.Role()
	}

	return m
}

func MapToChatMessage(m map[string]string) (ChatMessage, error) {
	switch ChatMessageType(m["type"]) {
	case ChatMessageTypeHuman:
		return NewHumanChatMessage(m["text"]), nil
	case ChatMessageTypeAI:
		return NewAIChatMessage(m["text"]), nil
	case ChatMessageTypeSystem:
		return NewSystemChatMessage(m["text"]), nil
	case ChatMessageTypeGeneric:
		return NewGenericChatMessage(m["text"], m["role"]), nil
	default:
		return nil, fmt.Errorf("unknown chat message type: %s", m["type"])
	}
}

type HumanChatMessage struct {
	text string
}

func NewHumanChatMessage(text string) *HumanChatMessage {
	return &HumanChatMessage{
		text: text,
	}
}

func (m HumanChatMessage) Type() ChatMessageType { return ChatMessageTypeHuman }
func (m HumanChatMessage) Text() string          { return m.text }

type AIChatMessage struct {
	text string
}

func NewAIChatMessage(text string) *AIChatMessage {
	return &AIChatMessage{
		text: text,
	}
}

func (m AIChatMessage) Type() ChatMessageType { return ChatMessageTypeAI }
func (m AIChatMessage) Text() string          { return m.text }

type SystemChatMessage struct {
	text string
}

func NewSystemChatMessage(text string) *SystemChatMessage {
	return &SystemChatMessage{
		text: text,
	}
}

func (m SystemChatMessage) Type() ChatMessageType { return ChatMessageTypeSystem }
func (m SystemChatMessage) Text() string          { return m.text }

type GenericChatMessage struct {
	text string
	role string
}

func NewGenericChatMessage(text, role string) *GenericChatMessage {
	return &GenericChatMessage{
		text: text,
		role: role,
	}
}

func (m GenericChatMessage) Type() ChatMessageType { return ChatMessageTypeGeneric }
func (m GenericChatMessage) Text() string          { return m.text }
func (m GenericChatMessage) Role() string          { return m.role }

type ChatMessages []ChatMessage

type StringifyChatMessagesOptions struct {
	HumanPrefix  string
	AIPrefix     string
	SystemPrefix string
}

func (cm ChatMessages) Format(optFns ...func(o *StringifyChatMessagesOptions)) (string, error) {
	opts := StringifyChatMessagesOptions{
		HumanPrefix:  "Human",
		AIPrefix:     "AI",
		SystemPrefix: "System",
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
		default:
			return "", fmt.Errorf("unknown chat message type: %s", message.Type())
		}

		result = append(result, fmt.Sprintf("%s: %s", role, message.Text()))
	}

	return strings.Join(result, "\n"), nil
}
