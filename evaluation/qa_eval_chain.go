package evaluation

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const qaEvalTemplate = `You are a teacher grading a quiz.
You are given a question, the student's answer, and the true answer, and are asked to score the student answer as either CORRECT or INCORRECT.

Example Format:
QUESTION: question here
STUDENT ANSWER: student's answer here
TRUE ANSWER: true answer here
GRADE: CORRECT or INCORRECT here

Grade the student answers based ONLY on their factual accuracy. Ignore differences in punctuation and phrasing between the student answer and true answer. It is OK if the student answer contains more information than the true answer, as long as it does not contain any conflicting statements. Begin! 

QUESTION: {{.query}}
STUDENT ANSWER: {{.result}}
TRUE ANSWER: {{.answer}}
GRADE:`

type QAEvalChainOptions struct {
	Prompt        schema.PromptTemplate
	QuestionKey   string
	AnswerKey     string
	PredictionKey string
}

// QAEvalChain is a LLM Chain specifically for evaluating question answering.
type QAEvalChain struct {
	llmChain      *chain.LLM
	questionKey   string
	answerKey     string
	predictionKey string
}

func NewQAEvalChain(model schema.Model, optFns ...func(o *QAEvalChainOptions)) (*QAEvalChain, error) {
	opts := QAEvalChainOptions{
		Prompt:        prompt.NewTemplate(qaEvalTemplate),
		QuestionKey:   "query",
		AnswerKey:     "answer",
		PredictionKey: "result",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	llmChain, err := chain.NewLLM(model, opts.Prompt)
	if err != nil {
		return nil, err
	}

	eval := &QAEvalChain{
		llmChain:      llmChain,
		questionKey:   opts.QuestionKey,
		answerKey:     opts.AnswerKey,
		predictionKey: opts.PredictionKey,
	}

	return eval, nil
}

func (eval *QAEvalChain) Evaluate(ctx context.Context, examples, predictions []map[string]string) ([]schema.ChainValues, error) {
	inputs := make([]schema.ChainValues, len(examples))

	for i, example := range examples {
		inputs[i] = schema.ChainValues{
			"query":  example[eval.questionKey],
			"answer": example[eval.answerKey],
			"result": predictions[i][eval.predictionKey],
		}
	}

	return golc.BatchCall(ctx, eval.llmChain, inputs)
}

func (eval *QAEvalChain) QuestionKey() string {
	return eval.questionKey
}

func (eval *QAEvalChain) AnswerKey() string {
	return eval.answerKey
}

func (eval *QAEvalChain) PredictionKey() string {
	return eval.predictionKey
}
