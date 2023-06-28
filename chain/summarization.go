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
	*schema.CallbackOptions
}

func NewStuffSummarization(llm schema.LLM, optFns ...func(o *StuffSummarizationOptions)) (*StuffDocuments, error) {
	opts := StuffSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuffPrompt := prompt.NewTemplate(stuffSummarizationTemplate)

	llmChain, err := NewLLM(llm, stuffPrompt, func(o *LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewStuffDocuments(llmChain, func(o *StuffDocumentsOptions) {
		o.CallbackOptions = opts.CallbackOptions
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
	*schema.CallbackOptions
}

func NewRefineSummarization(llm schema.LLM, optFns ...func(o *RefineSummarizationOptions)) (*RefineDocuments, error) {
	opts := RefineSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuffPrompt := prompt.NewTemplate(stuffSummarizationTemplate)

	llmChain, err := NewLLM(llm, stuffPrompt, func(o *LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	refinePrompt := prompt.NewTemplate(refineSummarizationTemplate)

	refineLLMChain, err := NewLLM(llm, refinePrompt, func(o *LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewRefineDocuments(llmChain, refineLLMChain, func(o *RefineDocumentsOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
}
