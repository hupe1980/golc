package rag

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultcondenseQuestionPromptTemplate = `Given the following conversation and a follow up question, rephrase the follow up question to be a standalone question, in its original language.

Chat History:
{{.history}}
Follow Up Input: {{.query}}
Standalone question:`

// Compile time check to ensure ConversationalRetrievalQA satisfies the Chain interface.
var _ schema.Chain = (*ConversationalRetrievalQA)(nil)

// ConversationalRetrievalQAOptions represents the options for the ConversationalRetrievalQA chain.
type ConversationalRetrievalQAOptions struct {
	*schema.CallbackOptions

	// Return the source documents
	ReturnSourceDocuments bool

	// Return the generated question
	ReturnGeneratedQuestion bool

	CondenseQuestionPrompt schema.PromptTemplate
	RetrievalQAPrompt      schema.PromptTemplate
	Memory                 schema.Memory
	InputKey               string
	OutputKey              string

	// If set, restricts the docs to return from store based on tokens, enforced only
	// for StuffDocumentsChain
	MaxTokenLimit uint
}

// ConversationalRetrievalQA is a chain implementation for conversational retrieval.
type ConversationalRetrievalQA struct {
	condenseQuestionChain *chain.LLM
	retrievalQAChain      *RetrievalQA
	opts                  ConversationalRetrievalQAOptions
}

// NewConversationalRetrievalQA creates a new instance of the ConversationalRetrievalQA chain.
func NewConversationalRetrievalQA(llm schema.LLM, retriever schema.Retriever, optFns ...func(o *ConversationalRetrievalQAOptions)) (*ConversationalRetrievalQA, error) {
	opts := ConversationalRetrievalQAOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ReturnSourceDocuments:   false,
		ReturnGeneratedQuestion: false,
		InputKey:                "question",
		OutputKey:               "answer",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Memory == nil {
		opts.Memory = memory.NewConversationBuffer(func(o *memory.ConversationBufferOptions) {
			o.OutputKey = opts.OutputKey
		})
	}

	if opts.CondenseQuestionPrompt == nil {
		opts.CondenseQuestionPrompt = prompt.NewTemplate(defaultcondenseQuestionPromptTemplate)
	}

	condenseQuestionChain, err := chain.NewLLM(llm, opts.CondenseQuestionPrompt)
	if err != nil {
		return nil, err
	}

	retrievalQAChain, err := NewRetrievalQA(llm, retriever, func(o *RetrievalQAOptions) {
		o.RetrievalQAPrompt = opts.RetrievalQAPrompt
		o.ReturnSourceDocuments = opts.ReturnSourceDocuments
		o.MaxTokenLimit = opts.MaxTokenLimit
		o.InputKey = opts.InputKey
	})
	if err != nil {
		return nil, err
	}

	return &ConversationalRetrievalQA{
		condenseQuestionChain: condenseQuestionChain,
		retrievalQAChain:      retrievalQAChain,
		opts:                  opts,
	}, nil
}

// Call executes the ConversationalRetrievalQA chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *ConversationalRetrievalQA) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	generatedQuestion, err := c.generateQuestion(ctx, inputs, opts)
	if err != nil {
		return nil, err
	}

	retrievalOutput, err := golc.Call(ctx, c.retrievalQAChain, schema.ChainValues{
		c.retrievalQAChain.InputKeys()[0]: generatedQuestion,
	}, func(co *golc.CallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	answer, err := retrievalOutput.GetString(c.retrievalQAChain.OutputKeys()[0])
	if err != nil {
		return nil, err
	}

	returns := schema.ChainValues{
		c.opts.OutputKey: answer,
	}

	if c.opts.ReturnSourceDocuments {
		returns["sourceDocuments"] = retrievalOutput["sourceDocuments"]
	}

	if c.opts.ReturnGeneratedQuestion {
		returns["generatedQuestion"] = generatedQuestion
	}

	return returns, nil
}

func (c *ConversationalRetrievalQA) generateQuestion(ctx context.Context, inputs schema.ChainValues, opts schema.CallOptions) (string, error) {
	if inputs["history"] == "" {
		return inputs.GetString(c.opts.InputKey)
	}

	output, err := golc.Call(ctx, c.condenseQuestionChain, inputs, func(co *golc.CallOptions) {
		co.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		co.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return "", err
	}

	return output.GetString(c.condenseQuestionChain.OutputKeys()[0])
}

// Memory returns the memory associated with the chain.
func (c *ConversationalRetrievalQA) Memory() schema.Memory {
	return c.opts.Memory
}

// Type returns the type of the chain.
func (c *ConversationalRetrievalQA) Type() string {
	return "ConversationalRetrievalQA"
}

// Verbose returns the verbosity setting of the chain.
func (c *ConversationalRetrievalQA) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *ConversationalRetrievalQA) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *ConversationalRetrievalQA) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *ConversationalRetrievalQA) OutputKeys() []string {
	outputKeys := []string{c.opts.OutputKey}
	if c.opts.ReturnSourceDocuments {
		outputKeys = append(outputKeys, "sourceDocuments")
	}

	if c.opts.ReturnGeneratedQuestion {
		outputKeys = append(outputKeys, "generatedQuestion")
	}

	return outputKeys
}
