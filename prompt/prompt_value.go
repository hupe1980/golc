// Package prompt provides utilities for managing and optimizing prompts.
package prompt

import (
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StringPromptValue satisfies the PromptValue interface.
var _ schema.PromptValue = (*StringPromptValue)(nil)

// StringPromptValue represents a string value that satisfies the PromptValue interface.
type StringPromptValue string

// String returns the string representation of the StringPromptValue.
func (v StringPromptValue) String() string {
	return string(v)
}

// Messages returns a ChatMessages slice containing a single HumanChatMessage with the string value.
func (v StringPromptValue) Messages() schema.ChatMessages {
	return schema.ChatMessages{
		schema.NewHumanChatMessage(string(v)),
	}
}

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
