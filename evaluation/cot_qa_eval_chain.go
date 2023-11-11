package evaluation

import (
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const cotQAEvalTemplate = `You are a teacher grading a quiz.
You are given a question, the context the question is about, and the student's answer. You are asked to score the student's answer as either CORRECT or INCORRECT, based on the context.
Write out in a step by step manner your reasoning to be sure that your conclusion is correct. Avoid simply stating the correct answer at the outset.

Example Format:
QUESTION: question here
CONTEXT: context the question is about here
STUDENT ANSWER: student's answer here
EXPLANATION: step by step reasoning here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. Ignore differences in punctuation and phrasing between the student answer and true answer. It is OK if the student answer contains more information than the true answer, as long as it does not contain any conflicting statements. Begin! 

QUESTION: {query}
CONTEXT: {context}
STUDENT ANSWER: {result}
EXPLANATION:`

type COTQAEvalChainOptions struct {
	Prompt        schema.PromptTemplate
	QuestionKey   string
	ContextKey    string
	PredictionKey string
}

// COTQAEvalChain is a LLM Chain specifically for evaluating QA using chain of thought reasoning.
type COTQAEvalChain struct {
	*ContextQAEvalChain
}

func NewCOTQAEvalChain(model schema.Model, optFns ...func(o *COTQAEvalChainOptions)) (*COTQAEvalChain, error) {
	opts := COTQAEvalChainOptions{
		QuestionKey:   "query",
		ContextKey:    "context",
		PredictionKey: "result",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Prompt == nil {
		opts.Prompt = prompt.NewTemplate(cotQAEvalTemplate)
	}

	contextQAEvalChain, err := NewContextQAEvalChain(model, func(o *ContextQAEvalChainOptions) {
		o.Prompt = opts.Prompt
		o.QuestionKey = opts.QuestionKey
		o.ContextKey = opts.ContextKey
		o.PredictionKey = opts.PredictionKey
	})
	if err != nil {
		return nil, err
	}

	eval := &COTQAEvalChain{}
	eval.ContextQAEvalChain = contextQAEvalChain

	return eval, nil
}
