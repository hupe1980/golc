package evaluation

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const contextQAEvalTemplate = `You are a teacher grading a quiz.
You are given a question, the context the question is about, and the student's answer. 
You are asked to score the student's answer as either CORRECT or INCORRECT, based on the context.

Example Format:
QUESTION: question here
CONTEXT: context the question is about here
STUDENT ANSWER: student's answer here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. 
Ignore differences in punctuation and phrasing between the student answer and true answer. 
It is OK if the student answer contains more information than the true answer, as long as 
it does not contain any conflicting statements. Begin! 

QUESTION: {{.query}}
CONTEXT: {{.context}}
STUDENT ANSWER: {{.result}}
GRADE:`

type ContextQAEvalChainOptions struct {
	Prompt        schema.PromptTemplate
	QuestionKey   string
	ContextKey    string
	PredictionKey string
}

// ConetxtQAEvalChain is a LLM Chain specifically for evaluating QA w/o GT based on context.
type ContextQAEvalChain struct {
	llmChain      *chain.LLM
	questionKey   string
	contextKey    string
	predictionKey string
}

func NewContextQAEvalChain(model schema.Model, optFns ...func(o *ContextQAEvalChainOptions)) (*ContextQAEvalChain, error) {
	opts := ContextQAEvalChainOptions{
		QuestionKey:   "query",
		ContextKey:    "context",
		PredictionKey: "result",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Prompt == nil {
		opts.Prompt = prompt.NewTemplate(contextQAEvalTemplate)
	}

	llmChain, err := chain.NewLLM(model, opts.Prompt)
	if err != nil {
		return nil, err
	}

	eval := &ContextQAEvalChain{
		llmChain:      llmChain,
		questionKey:   opts.QuestionKey,
		contextKey:    opts.ContextKey,
		predictionKey: opts.PredictionKey,
	}

	return eval, nil
}

func (eval *ContextQAEvalChain) Evaluate(ctx context.Context, examples, predictions []map[string]string) ([]schema.ChainValues, error) {
	inputs := make([]schema.ChainValues, len(examples))

	for i, example := range examples {
		inputs[i] = schema.ChainValues{
			"query":   example[eval.questionKey],
			"context": example[eval.contextKey],
			"result":  predictions[i][eval.predictionKey],
		}
	}

	return golc.BatchCall(ctx, eval.llmChain, inputs)
}
