package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultcondenseQuestionPromptTemplate = `Given the following conversation and a follow up question, rephrase the follow up question to be a standalone question, in its original language.

Chat History:
{{.history}}
Follow Up Input: {{.query}}
Standalone question:`

type ConversationalRetrievalOptions struct {
	*schema.CallbackOptions
	ReturnSourceDocuments   bool
	ReturnGeneratedQuestion bool
	CondenseQuestionPrompt  *prompt.Template
	Memory                  schema.Memory
	InputKey                string
	OutputKey               string
}

type ConversationalRetrieval struct {
	condenseQuestionChain *LLM
	retrievalQAChain      *RetrievalQA
	opts                  ConversationalRetrievalOptions
}

func NewConversationalRetrieval(llm schema.LLM, retriever schema.Retriever, optFns ...func(o *ConversationalRetrievalOptions)) (*ConversationalRetrieval, error) {
	opts := ConversationalRetrievalOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ReturnSourceDocuments:   false,
		ReturnGeneratedQuestion: false,
		InputKey:                "query",
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
		p, err := prompt.NewTemplate(defaultcondenseQuestionPromptTemplate)
		if err != nil {
			return nil, err
		}

		opts.CondenseQuestionPrompt = p
	}

	condenseQuestionChain, err := NewLLM(llm, opts.CondenseQuestionPrompt)
	if err != nil {
		return nil, err
	}

	retrievalQAChain, err := NewRetrievalQA(llm, retriever, func(o *RetrievalQAOptions) {
		o.ReturnSourceDocuments = opts.ReturnSourceDocuments
		o.InputKey = opts.InputKey
	})
	if err != nil {
		return nil, err
	}

	return &ConversationalRetrieval{
		condenseQuestionChain: condenseQuestionChain,
		retrievalQAChain:      retrievalQAChain,
		opts:                  opts,
	}, nil
}

func (c ConversationalRetrieval) Call(ctx context.Context, inputs schema.ChainValues) (schema.ChainValues, error) {
	output, err := golc.Call(ctx, c.condenseQuestionChain, inputs)
	if err != nil {
		return nil, err
	}

	generatedQuestion, ok := output[c.condenseQuestionChain.OutputKeys()[0]].(string)
	if !ok {
		return nil, fmt.Errorf("cannot convert generated question from output: %v", generatedQuestion)
	}

	retrievalOutput, err := golc.Call(ctx, c.retrievalQAChain, schema.ChainValues{
		"query": generatedQuestion,
	})
	if err != nil {
		return nil, err
	}

	answer, ok := retrievalOutput[c.retrievalQAChain.OutputKeys()[0]].(string)
	if !ok {
		return nil, fmt.Errorf("cannot convert answer from output: %v", generatedQuestion)
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

func (c ConversationalRetrieval) Memory() schema.Memory {
	return c.opts.Memory
}

func (c ConversationalRetrieval) Type() string {
	return "ConversationalRetrieval"
}

func (c ConversationalRetrieval) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

func (c ConversationalRetrieval) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c ConversationalRetrieval) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c ConversationalRetrieval) OutputKeys() []string {
	outputKeys := []string{c.opts.OutputKey}
	if c.opts.ReturnSourceDocuments {
		outputKeys = append(outputKeys, "sourceDocuments")
	}

	if c.opts.ReturnGeneratedQuestion {
		outputKeys = append(outputKeys, "generatedQuestion")
	}

	return outputKeys
}
