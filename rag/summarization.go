package rag

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const stuffSummarizationTemplate = `Write a concise summary of the following:


"{{.text}}"


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

	llmChain, err := chain.NewLLM(llm, stuffPrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.OutputKey = "outputText"
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
We have the opportunity to refine the existing summary (only if needed) with some more context below.
------------
{{.text}}
------------
Given the new context, refine the original summary
If the context isn't useful, return the original summary`

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

	llmChain, err := chain.NewLLM(llm, stuffPrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.OutputKey = "outputText"
	})
	if err != nil {
		return nil, err
	}

	refinePrompt := prompt.NewTemplate(refineSummarizationTemplate)

	refineLLMChain, err := chain.NewLLM(llm, refinePrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewRefineDocuments(llmChain, refineLLMChain, func(o *RefineDocumentsOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
}

const mapReduceSummarizationTemplate = `Write a concise summary of the following:


"{{.text}}"


CONCISE SUMMARY:`

type MapReduceSummarizationOptions struct {
	*schema.CallbackOptions
}

func NewMapReduceSummarization(llm schema.LLM, optFns ...func(o *MapReduceSummarizationOptions)) (*MapReduceDocuments, error) {
	opts := MapReduceSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	mapReducePrompt := prompt.NewTemplate(mapReduceSummarizationTemplate)

	mapChain, err := chain.NewLLM(llm, mapReducePrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	combineChain, err := NewStuffSummarization(llm)
	if err != nil {
		return nil, err
	}

	return NewMapReduceDocuments(mapChain, combineChain)
}
