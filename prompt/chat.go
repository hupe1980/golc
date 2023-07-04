package prompt

import (
	"fmt"

	"github.com/hupe1980/golc/schema"
)

type ChatPromptValue struct {
	messages schema.ChatMessages
}

func NewChatPromptValue(messages schema.ChatMessages) *ChatPromptValue {
	return &ChatPromptValue{
		messages: messages,
	}
}

func (v ChatPromptValue) String() (string, error) {
	return v.messages.Format()
}

func (v ChatPromptValue) Messages() schema.ChatMessages {
	return v.messages
}

type ChatTemplate interface {
	FormatPrompt(values map[string]any) (*ChatPromptValue, error)
	Format(values map[string]any) (schema.ChatMessages, error)
}

type chatTemplateWrapper struct {
	chatTemplates []ChatTemplate
}

func NewChatTemplateWrapper(chatTemplates ...ChatTemplate) ChatTemplate {
	return &chatTemplateWrapper{
		chatTemplates: chatTemplates,
	}
}

func (ct *chatTemplateWrapper) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

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

type chatTemplate struct {
	messageTemplates []MessageTemplate
}

func NewChatTemplate(messageTemplates []MessageTemplate) ChatTemplate {
	return &chatTemplate{
		messageTemplates: messageTemplates,
	}
}

func (ct *chatTemplate) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

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

type messagesPlaceholder struct {
	inputKey string
}

func NewMessagesPlaceholder(inputKey string) ChatTemplate {
	return &messagesPlaceholder{
		inputKey: inputKey,
	}
}

func (ct *messagesPlaceholder) FormatPrompt(values map[string]any) (*ChatPromptValue, error) {
	messages, err := ct.Format(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

func (ct *messagesPlaceholder) Format(values map[string]any) (schema.ChatMessages, error) {
	messages, ok := values[ct.inputKey].(schema.ChatMessages)
	if !ok {
		return nil, fmt.Errorf("cannot get list of messages for key %s", ct.inputKey)
	}

	return messages, nil
}

type MessageTemplate interface {
	Format(values map[string]any) (schema.ChatMessage, error)
}

type SystemMessageTemplate struct {
	prompt *Template
}

func NewSystemMessageTemplate(template string) *SystemMessageTemplate {
	return &SystemMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *SystemMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewSystemChatMessage(text), nil
}

type AIMessageTemplate struct {
	prompt *Template
}

func NewAIMessageTemplate(template string) *AIMessageTemplate {
	return &AIMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *AIMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewAIChatMessage(text), nil
}

type HumanMessageTemplate struct {
	prompt *Template
}

func NewHumanMessageTemplate(template string) *HumanMessageTemplate {
	return &HumanMessageTemplate{
		prompt: NewTemplate(template),
	}
}

func (pt *HumanMessageTemplate) Format(values map[string]any) (schema.ChatMessage, error) {
	text, err := pt.prompt.Format(values)
	if err != nil {
		return nil, err
	}

	return schema.NewHumanChatMessage(text), nil
}
