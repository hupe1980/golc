package rag

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const defaultStuffSummarizationTemplate = `Write a concise summary of the following:


"{{.text}}"


CONCISE SUMMARY:`

type StuffSummarizationOptions struct {
	*schema.CallbackOptions
	StuffSummarizationPrompt schema.PromptTemplate
}

func NewStuffSummarization(model schema.Model, optFns ...func(o *StuffSummarizationOptions)) (*StuffDocuments, error) {
	opts := StuffSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StuffSummarizationPrompt == nil {
		opts.StuffSummarizationPrompt = prompt.NewTemplate(defaultStuffSummarizationTemplate)
	}

	llmChain, err := chain.NewLLM(model, opts.StuffSummarizationPrompt, func(o *chain.LLMOptions) {
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

const defaultRefineSummarizationTemplate = `Your job is to produce a final summary
We have provided an existing summary up to a certain point: "{{.existingAnswer}}"
We have the opportunity to refine the existing summary (only if needed) with some more context below.
------------
{{.text}}
------------
Given the new context, refine the original summary
If the context isn't useful, return the original summary`

type RefineSummarizationOptions struct {
	*schema.CallbackOptions
	StuffSummarizationPrompt  schema.PromptTemplate
	RefineSummarizationPrompt schema.PromptTemplate
}

func NewRefineSummarization(model schema.Model, optFns ...func(o *RefineSummarizationOptions)) (*RefineDocuments, error) {
	opts := RefineSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StuffSummarizationPrompt == nil {
		opts.StuffSummarizationPrompt = prompt.NewTemplate(defaultStuffSummarizationTemplate)
	}

	llmChain, err := chain.NewLLM(model, opts.StuffSummarizationPrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.OutputKey = "outputText"
	})
	if err != nil {
		return nil, err
	}

	if opts.RefineSummarizationPrompt == nil {
		opts.RefineSummarizationPrompt = prompt.NewTemplate(defaultRefineSummarizationTemplate)
	}

	refineLLMChain, err := chain.NewLLM(model, opts.RefineSummarizationPrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	return NewRefineDocuments(llmChain, refineLLMChain, func(o *RefineDocumentsOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
}

const defaultMapReduceSummarizationTemplate = `Write a concise summary of the following:


"{{.text}}"


CONCISE SUMMARY:`

type MapReduceSummarizationOptions struct {
	*schema.CallbackOptions
	MapReduceSummarizationPrompt schema.PromptTemplate
	StuffSummarizationPrompt     schema.PromptTemplate
}

func NewMapReduceSummarization(model schema.Model, optFns ...func(o *MapReduceSummarizationOptions)) (*MapReduceDocuments, error) {
	opts := MapReduceSummarizationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.MapReduceSummarizationPrompt == nil {
		opts.MapReduceSummarizationPrompt = prompt.NewTemplate(defaultMapReduceSummarizationTemplate)
	}

	mapChain, err := chain.NewLLM(model, opts.MapReduceSummarizationPrompt, func(o *chain.LLMOptions) {
		o.CallbackOptions = opts.CallbackOptions
	})
	if err != nil {
		return nil, err
	}

	combineChain, err := NewStuffSummarization(model, func(o *StuffSummarizationOptions) {
		o.StuffSummarizationPrompt = opts.StuffSummarizationPrompt
	})
	if err != nil {
		return nil, err
	}

	return NewMapReduceDocuments(mapChain, combineChain)
}
