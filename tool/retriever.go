package tool

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/retriever"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Retriever satisfies the Tool interface.
var _ schema.Tool = (*Retriever)(nil)

// RetrieverOptions contains options for configuring the Retriever tool.
type RetrieverOptions struct {
	*schema.CallbackOptions
	DocumentSeparator string
}

// Retriever is a tool that utilizes a retriever to fetch documents based on a query.
type Retriever struct {
	retriever   schema.Retriever
	name        string
	description string
	opts        RetrieverOptions
}

// NewRetriever creates a new Retriever instance using the provided retriever, name, and description, along with optional configuration options.
func NewRetriever(retriever schema.Retriever, name, description string, optFns ...func(o *RetrieverOptions)) *Retriever {
	opts := RetrieverOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		DocumentSeparator: "\n\n",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Retriever{
		retriever:   retriever,
		name:        name,
		description: description,
		opts:        opts,
	}
}

// Name returns the name of the tool.
func (t *Retriever) Name() string {
	return t.name
}

// Description returns the description of the tool.
func (t *Retriever) Description() string {
	return t.description
}

// ArgsType returns the type of the input argument expected by the tool.
func (t *Retriever) ArgsType() reflect.Type {
	return reflect.TypeOf("") // string
}

// Run executes the tool with the given input and returns the output.
func (t *Retriever) Run(ctx context.Context, input any) (string, error) {
	query, ok := input.(string)
	if !ok {
		return "", errors.New("illegal input type")
	}

	docs, err := retriever.Run(ctx, t.retriever, query)
	if err != nil {
		return "", err
	}

	contents := make([]string, len(docs))
	for i, doc := range docs {
		contents[i] = doc.PageContent
	}

	return strings.Join(contents, t.opts.DocumentSeparator), nil
}

// Verbose returns the verbosity setting of the tool.
func (t *Retriever) Verbose() bool {
	return t.opts.Verbose
}

// Callbacks returns the registered callbacks of the tool.
func (t *Retriever) Callbacks() []schema.Callback {
	return t.opts.Callbacks
}
