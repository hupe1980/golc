package rag

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/retriever"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

const defaultRetrievalQAPromptTemplate = `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

{{.text}}

Question: {{.question}}
Helpful Answer:`

// Compile time check to ensure RetrievalQA satisfies the Chain interface.
var _ schema.Chain = (*RetrievalQA)(nil)

type RetrievalQAOptions struct {
	*schema.CallbackOptions
	RetrievalQAPrompt schema.PromptTemplate
	InputKey          string

	// Return the source documents
	ReturnSourceDocuments bool

	// If set, restricts the docs to return from store based on tokens, enforced only
	// for StuffDocumentsChain
	MaxTokenLimit uint
}

type RetrievalQA struct {
	stuffDocumentsChain *StuffDocuments
	retriever           schema.Retriever
	opts                RetrievalQAOptions
}

func NewRetrievalQA(model schema.LLM, retriever schema.Retriever, optFns ...func(o *RetrievalQAOptions)) (*RetrievalQA, error) {
	opts := RetrievalQAOptions{
		InputKey:              "question",
		ReturnSourceDocuments: false,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.RetrievalQAPrompt == nil {
		selector := prompt.ConditionalPromptSelector{
			DefaultPrompt: prompt.NewTemplate(defaultRetrievalQAPromptTemplate),
		}

		opts.RetrievalQAPrompt = selector.GetPrompt(model)
	}

	llmChain, err := chain.NewLLM(model, opts.RetrievalQAPrompt)
	if err != nil {
		return nil, err
	}

	stuffDocumentsChain, err := NewStuffDocuments(llmChain)
	if err != nil {
		return nil, err
	}

	return &RetrievalQA{
		stuffDocumentsChain: stuffDocumentsChain,
		retriever:           retriever,
		opts:                opts,
	}, nil
}

// Call executes the ConversationalRetrieval chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *RetrievalQA) Call(ctx context.Context, values schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	question, err := values.GetString(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	docs, err := c.getDocuments(ctx, question, opts)
	if err != nil {
		return nil, err
	}

	result, err := golc.Call(ctx, c.stuffDocumentsChain, schema.ChainValues{
		"question":                           question,
		c.stuffDocumentsChain.InputKeys()[0]: docs,
	}, func(co *golc.CallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	if c.opts.ReturnSourceDocuments {
		result["sourceDocuments"] = docs
	}

	return result, nil
}

func (c *RetrievalQA) getDocuments(ctx context.Context, query string, opts schema.CallOptions) ([]schema.Document, error) {
	docs, err := retriever.Run(ctx, c.retriever, query, func(o *retriever.Options) {
		o.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		o.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	numDocs := len(docs)

	if c.opts.MaxTokenLimit > 0 {
		tokens := make([]uint, len(docs))

		for i, doc := range docs {
			t, err := c.stuffDocumentsChain.llmChain.GetNumTokens(doc.PageContent)
			if err != nil {
				return nil, err
			}

			tokens[i] = t
		}

		tokenCount := util.SumInt(tokens[:numDocs])
		for tokenCount > c.opts.MaxTokenLimit {
			numDocs--
			tokenCount -= tokens[numDocs]
		}
	}

	return docs[:numDocs], nil
}

// Memory returns the memory associated with the chain.
func (c *RetrievalQA) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *RetrievalQA) Type() string {
	return "RetrievalQA"
}

// Verbose returns the verbosity setting of the chain.
func (c *RetrievalQA) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *RetrievalQA) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *RetrievalQA) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *RetrievalQA) OutputKeys() []string {
	return c.stuffDocumentsChain.OutputKeys()
}
