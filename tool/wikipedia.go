package tool

import (
	"context"
	"errors"
	"reflect"

	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Wikipedia satisfies the Tool interface.
var _ schema.Tool = (*Wikipedia)(nil)

type Wikipedia struct {
	client *integration.Wikipedia
}

func NewWikipedia(client *integration.Wikipedia) *Wikipedia {
	return &Wikipedia{
		client: client,
	}
}

// Name returns the name of the tool.
func (t *Wikipedia) Name() string {
	return "Wikipedia"
}

// Description returns the description of the tool.
func (t *Wikipedia) Description() string {
	return `A wrapper around Wikipedia.
Useful for when you need to answer general questions about 
people, places, companies, facts, historical events, or other subjects. 
Input should be a search query.`
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *Wikipedia) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *Wikipedia) Run(ctx context.Context, input any) (string, error) {
	query, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	return t.client.Run(ctx, query)
}

// Verbose returns the verbosity setting of the tool.
func (t *Wikipedia) Verbose() bool {
	return false
}

// Callbacks returns the registered callbacks of the tool.
func (t *Wikipedia) Callbacks() []schema.Callback {
	return nil
}
