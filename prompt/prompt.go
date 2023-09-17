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
