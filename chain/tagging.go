package chain

import (
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultTaggingTemplate = `Extract the desired information from the following passage.

Only extract the properties mentioned in the 'InformationExtraction' function.

Passage:
{input}`

// Compile time check to ensure Tagging satisfies the Chain interface.
var _ schema.Chain = (*Tagging)(nil)

// Tagging is a chain that uses structured output to perform tagging on a passage.
// It extracts the desired information from the given passage using a structured output model.
type Tagging struct {
	// StructuredOutput is the underlying structured output chain used for tagging.
	*StructuredOutput
}

// NewTagging creates a new Tagging chain with the provided chat model, structured output data, and optional options.
// It returns a Tagging chain or an error if the creation fails.
func NewTagging(chatModel schema.ChatModel, data any, optFns ...func(o *StructuredOutputOptions)) (*Tagging, error) {
	pt := prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewHumanMessageTemplate(defaultTaggingTemplate),
	})

	so, err := NewStructuredOutput(chatModel, pt, []OutputCandidate{{
		Name:        "InformationExtraction",
		Description: "Extracts the relevant information from the passage.",
		Data:        data,
	}}, optFns...)
	if err != nil {
		return nil, err
	}

	return &Tagging{
		StructuredOutput: so,
	}, nil
}

// Type returns the type of the chain.
func (c *Tagging) Type() string {
	return "Tagging"
}
