package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/util"
)

type StuffDocumentsOptions struct {
	callbackOptions
	InputKey             string
	DocumentVariableName string
	Separator            string
}

type StuffDocumentsChain struct {
	llmChain *LLMChain
	opts     StuffDocumentsOptions
}

func NewStuffDocumentsChain(llmChain *LLMChain, optFns ...func(o *StuffDocumentsOptions)) (*StuffDocumentsChain, error) {
	opts := StuffDocumentsOptions{
		InputKey:             "inputDocuments",
		DocumentVariableName: "context",
		Separator:            "\n\n",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	stuff := &StuffDocumentsChain{
		llmChain: llmChain,
		opts:     opts,
	}

	return stuff, nil
}

func (stuff *StuffDocumentsChain) Call(ctx context.Context, values golc.ChainValues) (golc.ChainValues, error) {
	cm := callback.NewManager(stuff.opts.Callbacks, stuff.opts.Verbose)

	if err := cm.OnChainStart("StuffDocumentsChain", &values); err != nil {
		return nil, err
	}

	input, ok := values[stuff.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, stuff.opts.InputKey)
	}

	docs, ok := input.([]golc.Document)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	contents := []string{}
	for _, doc := range docs {
		contents = append(contents, doc.PageContent)
	}

	inputValues := util.CopyMap(values)
	inputValues[stuff.opts.DocumentVariableName] = strings.Join(contents, stuff.opts.Separator)

	if err := cm.OnChainEnd(&golc.ChainValues{"outputs": inputValues}); err != nil {
		return nil, err
	}

	return stuff.llmChain.Call(ctx, inputValues)
}

// InputKeys returns the expected input keys.
func (stuff *StuffDocumentsChain) InputKeys() []string {
	return []string{stuff.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (stuff *StuffDocumentsChain) OutputKeys() []string {
	return stuff.llmChain.OutputKeys()
}
