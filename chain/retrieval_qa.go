package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type RetrievalQAOptions struct {
	InputKey              string
	ReturnSourceDocuments bool
}

type RetrievalQA struct {
	stuffDocumentsChain *StuffDocumentsChain
	retriever           schema.Retriever
	opts                RetrievalQAOptions
}

func NewRetrievalQA(chain *StuffDocumentsChain, retriever schema.Retriever) (*RetrievalQA, error) {
	opts := RetrievalQAOptions{
		InputKey:              "query",
		ReturnSourceDocuments: false,
	}

	qa := &RetrievalQA{
		stuffDocumentsChain: chain,
		retriever:           retriever,
		opts:                opts,
	}

	return qa, nil
}

func NewRetrievalQAFromLLM(llm schema.LLM, retriever schema.Retriever) (*RetrievalQA, error) {
	stuffPrompt, err := prompt.NewTemplate(defaultStuffQAPromptTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, stuffPrompt)
	if err != nil {
		return nil, err
	}

	stuffDocumentChain, err := NewStuffDocumentsChain(llmChain)
	if err != nil {
		return nil, err
	}

	return NewRetrievalQA(stuffDocumentChain, retriever)
}

func (qa *RetrievalQA) Call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[qa.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, qa.opts.InputKey)
	}

	query, ok := input.(string)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	docs, err := qa.retriever.GetRelevantDocuments(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := qa.stuffDocumentsChain.Call(ctx, map[string]any{
		"question":       query,
		"inputDocuments": docs,
	})
	if err != nil {
		return nil, err
	}

	if qa.opts.ReturnSourceDocuments {
		result["sourceDocuments"] = docs
	}

	return result, nil
}

// InputKeys returns the expected input keys.
func (qa *RetrievalQA) InputKeys() []string {
	return []string{qa.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (qa *RetrievalQA) OutputKeys() []string {
	return qa.stuffDocumentsChain.OutputKeys()
}
