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

func (t *Wikipedia) Name() string {
	return "Wikipedia"
}

func (t *Wikipedia) Description() string {
	return `A wrapper around Wikipedia.
Useful for when you need to answer general questions about 
people, places, companies, facts, historical events, or other subjects. 
Input should be a search query.`
}

func (t *Wikipedia) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

func (t *Wikipedia) Run(ctx context.Context, input any) (string, error) {
	query, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	return t.client.Run(ctx, query)
}

func (t *Wikipedia) Verbose() bool {
	return false
}

func (t *Wikipedia) Callbacks() []schema.Callback {
	return nil
}
