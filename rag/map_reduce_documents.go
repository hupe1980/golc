package rag

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure MapReduceDocuments satisfies the Chain interface.
var _ schema.Chain = (*MapReduceDocuments)(nil)

type MapReduceDocumentsOptions struct {
	*schema.CallbackOptions
	InputKey             string
	DocumentVariableName string
}

type MapReduceDocuments struct {
	mapChain     *chain.LLM
	combineChain *StuffDocuments
	opts         MapReduceDocumentsOptions
}

func NewMapReduceDocuments(mapChain *chain.LLM, combineChain *StuffDocuments, optFns ...func(o *MapReduceDocumentsOptions)) (*MapReduceDocuments, error) {
	opts := MapReduceDocumentsOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:             "inputDocuments",
		DocumentVariableName: "text",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &MapReduceDocuments{
		mapChain:     mapChain,
		combineChain: combineChain,
		opts:         opts,
	}, nil
}

// Call executes the MapReduceDocuments chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *MapReduceDocuments) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
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

	rest := schema.ChainValues(util.OmitByKeys(inputs, []string{c.opts.InputKey}))

	batchInputs := make([]schema.ChainValues, len(docs))

	for i, d := range docs {
		batchInput := rest.Clone()
		batchInput[c.opts.DocumentVariableName] = d.PageContent
		batchInputs[i] = batchInput
	}

	mapResults, err := golc.BatchCall(ctx, c.mapChain, batchInputs, func(co *golc.BatchCallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	combineDocs := make([]schema.Document, len(docs))

	for i, d := range docs {
		mapResult, err := mapResults[i].GetString(c.mapChain.OutputKeys()[0])
		if err != nil {
			return nil, err
		}

		combineDocs[i] = schema.Document{
			PageContent: mapResult,
			Metadata:    d.Metadata,
		}
	}

	combineInputs := rest.Clone()
	combineInputs[c.combineChain.InputKeys()[0]] = combineDocs

	return golc.Call(ctx, c.combineChain, combineInputs, func(co *golc.CallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
}

// Memory returns the memory associated with the chain.
func (c *MapReduceDocuments) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *MapReduceDocuments) Type() string {
	return "MapReduceDocuments"
}

// Verbose returns the verbosity setting of the chain.
func (c *MapReduceDocuments) Verbose() bool {
	return c.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *MapReduceDocuments) Callbacks() []schema.Callback {
	return c.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (c *MapReduceDocuments) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *MapReduceDocuments) OutputKeys() []string {
	return c.combineChain.OutputKeys()
}
