package prompt

import (
	"fmt"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// ChatTemplate represents a chat  template.
type ChatTemplate interface {
	schema.PromptTemplate
	FormatMessages(values map[string]any) (schema.ChatMessages, error)
}

// Compile time check to ensure chatTemplateWrapper satisfies the PromptTemplate interface.
var _ schema.PromptTemplate = (*chatTemplateWrapper)(nil)

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

func (ct *chatTemplateWrapper) Format(values map[string]any) (string, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return "", err
	}

	return messages.Format()
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *chatTemplateWrapper) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// FormatMessages formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *chatTemplateWrapper) FormatMessages(values map[string]any) (schema.ChatMessages, error) {
	fullMeessages := schema.ChatMessages{}

	for _, t := range ct.chatTemplates {
		messages, err := t.FormatMessages(values)
		if err != nil {
			return nil, err
		}

		fullMeessages = append(fullMeessages, messages...)
	}

	return fullMeessages, nil
}

// InputVariables returns a list of input variables used by all the wrapped ChatTemplates.
func (ct *chatTemplateWrapper) InputVariables() []string {
	inputVariables := make([]string, 0)
	for _, ct := range ct.chatTemplates {
		inputVariables = append(inputVariables, ct.InputVariables()...)
	}

	return util.Uniq(inputVariables)
}

// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
func (ct *chatTemplateWrapper) OutputParser() (schema.OutputParser[any], bool) {
	return nil, false
}

// Compile time check to ensure chatTemplate satisfies the PromptTemplate interface.
var _ schema.PromptTemplate = (*chatTemplate)(nil)

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

func (ct *chatTemplate) Format(values map[string]any) (string, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return "", err
	}

	return messages.Format()
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *chatTemplate) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// FormatMessages formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *chatTemplate) FormatMessages(values map[string]any) (schema.ChatMessages, error) {
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

// InputVariables returns a list of input variables used by the message templates.
func (ct *chatTemplate) InputVariables() []string {
	inputVariables := make([]string, 0)
	for _, mt := range ct.messageTemplates {
		inputVariables = append(inputVariables, mt.InputVariables()...)
	}

	return util.Uniq(inputVariables)
}

// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
func (ct *chatTemplate) OutputParser() (schema.OutputParser[any], bool) {
	return nil, false
}

// Compile time check to ensure messagesPlaceholder satisfies the PromptTemplate interface.
var _ schema.PromptTemplate = (*messagesPlaceholder)(nil)

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

func (ct *messagesPlaceholder) Format(values map[string]any) (string, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return "", err
	}

	return messages.Format()
}

// FormatPrompt formats the prompt using the provided values and returns a ChatPromptValue.
func (ct *messagesPlaceholder) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
	messages, err := ct.FormatMessages(values)
	if err != nil {
		return nil, err
	}

	return NewChatPromptValue(messages), nil
}

// FormatMessages formats the messages using the provided values and returns the resulting ChatMessages.
func (ct *messagesPlaceholder) FormatMessages(values map[string]any) (schema.ChatMessages, error) {
	messages, ok := values[ct.inputKey].(schema.ChatMessages)
	if !ok {
		return nil, fmt.Errorf("cannot get list of messages for key %s", ct.inputKey)
	}

	return messages, nil
}

// InputVariables returns an empty list for the messagesPlaceholder since it doesn't use input variables.
func (ct *messagesPlaceholder) InputVariables() []string {
	return []string{}
}

// OutputParser returns the output parser function and a boolean indicating if an output parser is defined.
func (ct *messagesPlaceholder) OutputParser() (schema.OutputParser[any], bool) {
	return nil, false
}

// MessageTemplate represents a chat message template.
type MessageTemplate interface {
	Format(values map[string]any) (schema.ChatMessage, error)
	FormatPrompt(values map[string]any) (schema.PromptValue, error)
	InputVariables() []string
}

// Compile time check to ensure SystemMessageTemplate satisfies the MessageTemplate interface.
var _ MessageTemplate = (*SystemMessageTemplate)(nil)

// Compile time check to ensure AIMessageTemplate satisfies the MessageTemplate interface.
var _ MessageTemplate = (*AIMessageTemplate)(nil)

// Compile time check to ensure HumanMessageTemplate satisfies the MessageTemplate interface.
var _ MessageTemplate = (*HumanMessageTemplate)(nil)

type messageTemplate struct {
	MessageTemplate
}

func (mt *messageTemplate) FormatPrompt(values map[string]any) (schema.PromptValue, error) {
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
func NewSystemMessageTemplate(template string, optFns ...func(o *TemplateOptions)) *SystemMessageTemplate {
	opts := DefaultTemplateOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	mt := &SystemMessageTemplate{
		prompt: NewTemplate(template, func(o *TemplateOptions) {
			*o = opts
		}),
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

// InputVariables returns the input variables used in the system message template.
func (pt *SystemMessageTemplate) InputVariables() []string {
	return pt.prompt.InputVariables()
}

// AIMessageTemplate represents an AI message template.
type AIMessageTemplate struct {
	messageTemplate
	prompt *Template
}

// NewAIMessageTemplate creates a new AIMessageTemplate with the given template.
func NewAIMessageTemplate(template string, optFns ...func(o *TemplateOptions)) *AIMessageTemplate {
	opts := DefaultTemplateOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	mt := &AIMessageTemplate{
		prompt: NewTemplate(template, func(o *TemplateOptions) {
			*o = opts
		}),
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

// InputVariables returns the input variables used in the AI message template.
func (pt *AIMessageTemplate) InputVariables() []string {
	return pt.prompt.InputVariables()
}

// HumanMessageTemplate represents a human message template.
type HumanMessageTemplate struct {
	messageTemplate
	prompt *Template
}

// NewHumanMessageTemplate creates a new HumanMessageTemplate with the given template.
func NewHumanMessageTemplate(template string, optFns ...func(o *TemplateOptions)) *HumanMessageTemplate {
	opts := DefaultTemplateOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	mt := &HumanMessageTemplate{
		prompt: NewTemplate(template, func(o *TemplateOptions) {
			*o = opts
		}),
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

// InputVariables returns the input variables used in the human message template.
func (pt *HumanMessageTemplate) InputVariables() []string {
	return pt.prompt.InputVariables()
}
