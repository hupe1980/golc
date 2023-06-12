package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/util"
)

type StuffDocumentsOptions struct {
	InputKey             string
	DocumentVariableName string
	Separator            string
}

type StuffDocumentsChain struct {
	*Chain
	llmChain *LLMChain
	opts     StuffDocumentsOptions
}

func NewStuffDocumentsChain(llmChain *LLMChain) (*StuffDocumentsChain, error) {
	opts := StuffDocumentsOptions{
		InputKey:             "inputDocuments",
		DocumentVariableName: "context",
		Separator:            "\n\n",
	}

	stuff := &StuffDocumentsChain{
		llmChain: llmChain,
		opts:     opts,
	}

	stuff.Chain = NewChain(stuff.call)

	return stuff, nil
}

func (stuff *StuffDocumentsChain) call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	input, ok := values[stuff.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, stuff.opts.InputKey)
	}

	docs, ok := input.([]golc.Document)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	var sb strings.Builder
	for _, doc := range docs {
		sb.WriteString(doc.PageContent)
		sb.WriteString(stuff.opts.Separator)
	}

	inputValues := util.CopyMap(values)
	inputValues[stuff.opts.DocumentVariableName] = sb.String()

	return stuff.llmChain.Call(ctx, inputValues)
}
