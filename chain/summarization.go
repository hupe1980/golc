package chain

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const stuffSummarizationTemplate = `Write a concise summary of the following:


"{{.context}}"


CONCISE SUMMARY:`

type StuffSummarizationOptions struct {
	*callbackOptions
}

func NewStuffSummarization(llm schema.LLM, optFns ...func(o *StuffSummarizationOptions)) (*StuffDocuments, error) {
	opts := StuffSummarizationOptions{
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuffPrompt, err := prompt.NewTemplate(stuffSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, stuffPrompt, func(o *LLMChainOptions) {
		o.callbackOptions = opts.callbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewStuffDocuments(llmChain, func(o *StuffDocumentsOptions) {
		o.callbackOptions = opts.callbackOptions
	})
}

const refineSummarizationTemplate = `Your job is to produce a final summary
We have provided an existing summary up to a certain point: "{{.existingAnswer}}"
We have the opportunity to refine the existing summary
(only if needed) with some more context below.
------------
"{{.text}}"
------------

Given the new context, refine the original summary
If the context isn't useful, return the original summary.

REFINED SUMMARY:`

type RefineSummarizationOptions struct {
	*callbackOptions
}

func NewRefineSummarization(llm schema.LLM, optFns ...func(o *RefineSummarizationOptions)) (*RefineDocuments, error) {
	opts := RefineSummarizationOptions{
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuffPrompt, err := prompt.NewTemplate(stuffSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, stuffPrompt, func(o *LLMChainOptions) {
		o.callbackOptions = opts.callbackOptions
	})
	if err != nil {
		return nil, err
	}

	refinePrompt, err := prompt.NewTemplate(refineSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	refineLLMChain, err := NewLLMChain(llm, refinePrompt, func(o *LLMChainOptions) {
		o.callbackOptions = opts.callbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewRefineDocuments(llmChain, refineLLMChain, func(o *RefineDocumentsOptions) {
		o.callbackOptions = opts.callbackOptions
	})
}
