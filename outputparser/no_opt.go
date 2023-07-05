package outputparser

import (
	"fmt"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure NoOpts satisfies the OutputParser interface.
var _ schema.OutputParser[any] = (*NoOpt)(nil)

type NoOpt struct{}

func NewNoOpt() *NoOpt {
	return &NoOpt{}
}

func (p *NoOpt) ParseResult(result schema.Generation) (any, error) {
	return result.Text, nil
}

func (p *NoOpt) Parse(text string) (any, error) {
	return text, nil
}

func (p *NoOpt) ParseWithPrompt(text string, prompt schema.PromptValue) (any, error) {
	return p.Parse(text)
}

func (p *NoOpt) GetFormatInstructions() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (p *NoOpt) Type() string {
	return "no-opt"
}
