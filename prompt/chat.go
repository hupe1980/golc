package prompt

import (
	"fmt"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure ChatPromptValue satisfies the PromptValue interface.
var _ schema.PromptValue = (*ChatPromptValue)(nil)

// ChatPromptValue represents a chat prompt value containing chat messages.
type ChatPromptValue struct {
	messages schema.ChatMessages
}

// NewChatPromptValue creates a new ChatPromptValue with the given chat messages.
func NewChatPromptValue(messages schema.ChatMessages) *ChatPromptValue {
	return &ChatPromptValue{
		messages: messages,
	}
}

// String returns a string representation of the ChatPromptValue.
func (v ChatPromptValue) String() string {
	pv, err := v.messages.Format()
	if err != nil {
		panic(err)
	}

	return pv
}

// Messages returns the chat messages contained in the ChatPromptValue.
func (v ChatPromptValue) Messages() schema.ChatMessages {
	return v.messages
}

// ChatTemplate represents a chat  template.
type ChatTemplate interface {
	FormatPrompt(values map[string]any) (*ChatPromptValue, error)
	Format(values map[string]any) (schema.ChatMessages, error)
}

// chatTemplateWrapper wraps multiple ChatTemplates and provides combined formatting.
type chatTemplateWrapper struct {
	chatTemplates []ChatTemplate
}

// NewChatTemplateWrapper creates a new ChatTemplate that wraps multiple ChatTemplates.
func NewChatTemplateWrapper(chatTemplates ...ChatTemplate) ChatTemplate {
	return &chatTemplateWrapper{
		chatTemplates: chatTemplates,
	}
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *chatTemplateWrapper) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// Format formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *chatTemplateWrapper) Format(values map[string]any) (schema.ChatMessages, error) {
	fullMeessages := schema.ChatMessages{}

	for _, t := range ct.chatTemplates {
		messages, err := t.Format(values)
		if err != nil {
			return nil, err
		}

		fullMeessages = append(fullMeessages, messages...)
	}

	return fullMeessages, nil
}

// chatTemplate represents a chat message template.
type chatTemplate struct {
	messageTemplates []MessageTemplate
}

// NewChatTemplate creates a new ChatTemplate with the given message templates.
func NewChatTemplate(messageTemplates []MessageTemplate) ChatTemplate {
	return &chatTemplate{
		messageTemplates: messageTemplates,
	}
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *chatTemplate) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// Format formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *chatTemplate) Format(values map[string]any) (schema.ChatMessages, error) {
	messages := make(schema.ChatMessages, len(ct.messageTemplates))

	for i, t := range ct.messageTemplates {
		msg, err := t.Format(values)
		if err != nil {
			return nil, err
		}

		messages[i] = msg
	}

	return messages, nil
}

// messagesPlaceholder represents a placeholder for chat messages.
type messagesPlaceholder struct {
	inputKey string
}

// NewMessagesPlaceholder creates a new ChatTemplate placeholder for chat messages.
func NewMessagesPlaceholder(inputKey string) ChatTemplate {
	return &messagesPlaceholder{
		inputKey: inputKey,
	}
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *messagesPlaceholder) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// Format formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *messagesPlaceholder) Format(values map[string]any) (schema.ChatMessages, error) {
	messages, ok := values[ct.inputKey].(schema.ChatMessages)
	if !ok {
		return nil, fmt.Errorf("cannot get list of messages for key %s", ct.inputKey)
	}

	return messages, nil
}

// MessageTemplate represents a chat message template.
type MessageTemplate interface {
	Format(values map[string]any) (schema.ChatMessage, error)
	FormatPrompt(values map[string]any) (*ChatPromptValue, error)
}

type messageTemplate struct {
	MessageTemplate
}

func (mt *messageTemplate) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	message, err := mt.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(schema.ChatMessages{message}), nil
}

// SystemMessageTemplate represents a system message template.
type SystemMessageTemplate struct {
	messageTemplate
	prompt *Template
}

// NewSystemMessageTemplate creates a new SystemMessageTemplate with the given template.
func NewSystemMessageTemplate(template string) *SystemMessageTemplate {
	mt := &SystemMessageTemplate{
		prompt: NewTemplate(template),
	}

	mt.messageTemplate = messageTemplate{mt}

	return mt
}

// Format formats the message using the provided values and returns a SystemChatMessage.
func (pt *SystemMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewSystemChatMessage(text), nil
}

// AIMessageTemplate represents an AI message template.
type AIMessageTemplate struct {
	messageTemplate
	prompt *Template
}

// NewAIMessageTemplate creates a new AIMessageTemplate with the given template.
func NewAIMessageTemplate(template string) *AIMessageTemplate {
	mt := &AIMessageTemplate{
		prompt: NewTemplate(template),
	}

	mt.messageTemplate = messageTemplate{mt}

	return mt
}

// Format formats the message using the provided values and returns an AIChatMessage.
func (pt *AIMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewAIChatMessage(text), nil
}

// HumanMessageTemplate represents a human message template.
type HumanMessageTemplate struct {
	messageTemplate
	prompt *Template
}

// NewHumanMessageTemplate creates a new HumanMessageTemplate with the given template.
func NewHumanMessageTemplate(template string) *HumanMessageTemplate {
	mt := &HumanMessageTemplate{
		prompt: NewTemplate(template),
	}

	mt.messageTemplate = messageTemplate{mt}

	return mt
}

// Format formats the message using the provided values and returns a HumanChatMessage.
func (pt *HumanMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewHumanChatMessage(text), nil
}
