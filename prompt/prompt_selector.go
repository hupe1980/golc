package prompt

import (
	"github.com/hupe1980/golc/schema"
)

// PromptSelector is an interface for selecting prompts based on a model.
type PromptSelector interface {
	GetPrompt(model schema.Model) schema.PromptTemplate
}

// Compile time check to ensure ConditionalPromptSelector satisfies the PromptSelector interface.
var _ PromptSelector = (*ConditionalPromptSelector)(nil)

// ConditionalFunc represents a function that evaluates a condition based on a model.
type ConditionalFunc func(model schema.Model) bool

// Conditional represents a conditional prompt configuration.
type Conditional struct {
	Condition ConditionalFunc
	Prompt    schema.PromptTemplate
}

// ConditionalPromptSelector is a prompt selector that selects prompts based on conditions.
type ConditionalPromptSelector struct {
	DefaultPrompt schema.PromptTemplate
	Conditionals  []Conditional
}

// GetPrompt selects a prompt template based on the provided model.
// It evaluates the conditions in order and returns the prompt associated with the first matching condition,
// or returns the default prompt if no condition is met.
func (cps *ConditionalPromptSelector) GetPrompt(model schema.Model) schema.PromptTemplate {
	for _, conditional := range cps.Conditionals {
		if conditional.Condition(model) {
			return conditional.Prompt
		}
	}

	return cps.DefaultPrompt
}

// IsLLM checks if the given model is of type schema.LLM.
func IsLLM(model schema.Model) bool {
	_, ok := model.(schema.LLM)
	return ok
}

// IsChatModel checks if the given model is of type schema.ChatModel.
func IsChatModel(model schema.Model) bool {
	_, ok := model.(schema.ChatModel)
	return ok
}
