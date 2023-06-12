package chain

import (
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
)

const stuffSummarizationTemplate = `Write a concise summary of the following:


"{{.context}}"


CONCISE SUMMARY:`

func NewStuffSummarizationChain(llm golc.LLM) (*StuffDocumentsChain, error) {
	stuffPrompt, err := prompt.NewTemplate(stuffSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, stuffPrompt)
	if err != nil {
		return nil, err
	}

	return NewStuffDocumentsChain(llmChain)
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

func NewRefineSummarizationChain(llm golc.LLM) (*RefineDocumentsChain, error) {
	stuffPrompt, err := prompt.NewTemplate(stuffSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	llmChain, err := NewLLMChain(llm, stuffPrompt)
	if err != nil {
		return nil, err
	}

	refinePrompt, err := prompt.NewTemplate(refineSummarizationTemplate)
	if err != nil {
		return nil, err
	}

	refineLLMChain, err := NewLLMChain(llm, refinePrompt)
	if err != nil {
		return nil, err
	}

	return NewRefineDocumentsChain(llmChain, refineLLMChain)
}
