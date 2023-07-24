package outputparser

import (
	"errors"
	"strings"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure CommaSeparatedList satisfies the OutputParser interface.
var _ schema.OutputParser[any] = (*CommaSeparatedList)(nil)

// CommaSeparatedList is an implementation of the OutputParser interface that parses
// a comma-separated list of values from the output text.
type CommaSeparatedList struct{}

// NewCommaSeparatedList creates a new instance of the CommaSeparatedList parser.
func NewCommaSeparatedList() CommaSeparatedList {
	return CommaSeparatedList{}
}

// ParseResult parses the result from code generation into a comma-separated list of values.
// It implements the ParseResult method of the OutputParser interface.
func (p *CommaSeparatedList) ParseResult(result schema.Generation) (any, error) {
	return p.Parse(result.Text)
}

// Parse parses the input text as a comma-separated list of values and returns
// them as a slice of strings. The input text is expected to be in the format of
// a comma-separated list, e.g., "foo, bar, baz". Leading and trailing spaces in
// each value will be removed during parsing.
//
// If the input text is empty or contains only spaces, it will return an error
// with the message "no value to parse".
//
// It implements the Parse method of the OutputParser interface.
func (p *CommaSeparatedList) Parse(text string) (any, error) {
	input := strings.TrimSpace(text)
	if input == "" {
		return nil, errors.New("no value to parse")
	}

	values := strings.Split(input, ",")
	for i := 0; i < len(values); i++ {
		values[i] = strings.TrimSpace(values[i])
	}

	return values, nil
}

// ParseWithPrompt parses a comma-separated list of values from the provided text and prompt.
// It implements the ParseWithPrompt method of the OutputParser interface.
func (p *CommaSeparatedList) ParseWithPrompt(text string, prompt schema.PromptValue) (any, error) {
	return p.Parse(text)
}

// GetFormatInstructions returns the format instructions for using the CommaSeparatedList parser.
// It implements the GetFormatInstructions method of the OutputParser interface.
func (p *CommaSeparatedList) GetFormatInstructions() string {
	return "Your response should be a list of comma-separated values, e.g.: `foo, bar, baz`"
}

// Type returns the type of the output parser, which is "comma_separated_list".
func (p *CommaSeparatedList) Type() string {
	return "comma_separated_list"
}
