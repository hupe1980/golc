package outputparser

import (
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure NoOpts satisfies the OutputParser interface.
var _ schema.OutputParser[any] = (*NoOpt)(nil)

// NoOpt represents a simple parser that returns the output text without any additional processing.
type NoOpt struct{}

// NewNoOpt creates a new instance of NoOpt parser.
func NewNoOpt() *NoOpt {
	return &NoOpt{}
}

// ParseResult returns the generation text as the parsed result without any modifications.
func (p *NoOpt) ParseResult(result schema.Generation) (any, error) {
	return result.Text, nil
}

// Parse simply returns the input text as the parsed result without any modifications.
func (p *NoOpt) Parse(text string) (any, error) {
	return text, nil
}

// ParseWithPrompt is not used for this parser, so it simply calls Parse.
func (p *NoOpt) ParseWithPrompt(text string, prompt schema.PromptValue) (any, error) {
	return p.Parse(text)
}

// GetFormatInstructions returns an empty string as there are no specific format instructions for this parser.
func (p *NoOpt) GetFormatInstructions() string {
	return ""
}

// Type returns the type identifier of the parser, which is "no_opt".
func (p *NoOpt) Type() string {
	return "no_opt"
}
