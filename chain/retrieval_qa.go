package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type RetrievalQAOptions struct {
	*schema.CallbackOptions
	InputKey              string
	ReturnSourceDocuments bool
}

type RetrievalQA struct {
	stuffDocumentsChain *StuffDocuments
	retriever           schema.Retriever
	opts                RetrievalQAOptions
}

func NewRetrievalQA(stuffDocumentsChain *StuffDocuments, retriever schema.Retriever, optFns ...func(o *RetrievalQAOptions)) (*RetrievalQA, error) {
	opts := RetrievalQAOptions{
		InputKey:              "query",
		ReturnSourceDocuments: false,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &RetrievalQA{
		stuffDocumentsChain: stuffDocumentsChain,
		retriever:           retriever,
		opts:                opts,
	}, nil
}

func NewRetrievalQAFromLLM(llm schema.LLM, retriever schema.Retriever, optFns ...func(o *RetrievalQAOptions)) (*RetrievalQA, error) {
	stuffPrompt, err := prompt.NewTemplate(defaultStuffQAPromptTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLM(llm, stuffPrompt)
	if err != nil {
		return nil, err
	}

	stuffDocumentChain, err := NewStuffDocuments(llmChain)
	if err != nil {
		return nil, err
	}

	return NewRetrievalQA(stuffDocumentChain, retriever, optFns...)
}

func (c *RetrievalQA) Call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	query, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	docs, err := c.retriever.GetRelevantDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := golc.Call(ctx, c.stuffDocumentsChain, map[string]any{
		"question":       query,
		"inputDocuments": docs,
	})
	if err != nil {
		return nil, err
	}

	if c.opts.ReturnSourceDocuments {
		result["sourceDocuments"] = docs
	}

	return result, nil
}

func (c *RetrievalQA) Memory() schema.Memory {
	return nil
}

func (c *RetrievalQA) Type() string {
	return "RetrievalQA"
}

func (c *RetrievalQA) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

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
