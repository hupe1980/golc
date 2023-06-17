package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type StuffDocumentsOptions struct {
	*callbackOptions
	InputKey             string
	DocumentVariableName string
	Separator            string
}

type StuffDocumentsChain struct {
	*baseChain
	llmChain *LLMChain
	opts     StuffDocumentsOptions
}

func NewStuffDocumentsChain(llmChain *LLMChain, optFns ...func(o *StuffDocumentsOptions)) (*StuffDocumentsChain, error) {
	opts := StuffDocumentsOptions{
		InputKey:             "inputDocuments",
		DocumentVariableName: "context",
		Separator:            "\n\n",
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuff := &StuffDocumentsChain{
		llmChain: llmChain,
		opts:     opts,
	}

	stuff.baseChain = &baseChain{
		chainName:       "StuffDocumentsChain",
		callFunc:        stuff.call,
		inputKeys:       []string{opts.InputKey},
		outputKeys:      llmChain.OutputKeys(),
		callbackOptions: opts.callbackOptions,
	}

	return stuff, nil
}

func (stuff *StuffDocumentsChain) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[stuff.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, stuff.opts.InputKey)
	}

	docs, ok := input.([]schema.Document)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	contents := []string{}
	for _, doc := range docs {
		contents = append(contents, doc.PageContent)
	}

	inputValues := util.CopyMap(values)
	inputValues[stuff.opts.DocumentVariableName] = strings.Join(contents, stuff.opts.Separator)

	return stuff.llmChain.Call(ctx, inputValues)
}
