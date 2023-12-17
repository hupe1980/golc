package rag

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StuffDocuments satisfies the Chain interface.
var _ schema.Chain = (*StuffDocuments)(nil)

type StuffDocumentsOptions struct {
	*schema.CallbackOptions
	InputKey             string
	DocumentVariableName string
	DocumentSeparator    string
}

type StuffDocuments struct {
	llmChain *chain.LLM
	opts     StuffDocumentsOptions
}

func NewStuffDocuments(llmChain *chain.LLM, optFns ...func(o *StuffDocumentsOptions)) (*StuffDocuments, error) {
	opts := StuffDocumentsOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:             "inputDocuments",
		DocumentVariableName: "text",
		DocumentSeparator:    "\n\n",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &StuffDocuments{
		llmChain: llmChain,
		opts:     opts,
	}, nil
}

// Call executes the StuffDocuments chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *StuffDocuments) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	docs, err := inputs.GetDocuments(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	contents := make([]string, len(docs))
	for i, doc := range docs {
		contents[i] = doc.PageContent
	}

	rest := schema.ChainValues(util.OmitByKeys(inputs, []string{c.opts.InputKey}))

	rest[c.opts.DocumentVariableName] = strings.Join(contents, c.opts.DocumentSeparator)

	output, err := golc.SimpleCall(ctx, c.llmChain, rest, func(co *golc.SimpleCallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.llmChain.OutputKeys()[0]: strings.TrimSpace(output),
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *StuffDocuments) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *StuffDocuments) Type() string {
	return "StuffDocuments"
}

// Verbose returns the verbosity setting of the chain.
func (c *StuffDocuments) Verbose() bool {
	return c.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *StuffDocuments) Callbacks() []schema.Callback {
	return c.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (c *StuffDocuments) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *StuffDocuments) OutputKeys() []string {
	return c.llmChain.OutputKeys()
}
