package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type RefineDocumentsOptions struct {
	InputKey             string
	DocumentVariableName string
	InitialResponseName  string
	DocumentPrompt       *prompt.Template
	OutputKey            string
}

type RefineDocumentsChain struct {
	*chain
	llmChain       *LLMChain
	refineLLMChain *LLMChain
	opts           RefineDocumentsOptions
}

func NewRefineDocumentsChain(llmChain *LLMChain, refineLLMChain *LLMChain) (*RefineDocumentsChain, error) {
	opts := RefineDocumentsOptions{
		InputKey:             "inputDocuments",
		DocumentVariableName: "context",
		InitialResponseName:  "existingAnswer",
		OutputKey:            "text",
	}

	if opts.DocumentPrompt == nil {
		p, err := prompt.NewTemplate("{{.pageContent}}")
		if err != nil {
			return nil, err
		}

		opts.DocumentPrompt = p
	}

	refine := &RefineDocumentsChain{
		llmChain:       llmChain,
		refineLLMChain: refineLLMChain,
		opts:           opts,
	}

	refine.chain = newChain(refine.call, []string{opts.InputKey}, llmChain.OutputKeys())

	return refine, nil
}

func (refine *RefineDocumentsChain) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[refine.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, refine.opts.InputKey)
	}

	docs, ok := input.([]schema.Document)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("%w: documents slice has no elements", ErrInvalidInputValues)
	}

	rest := util.OmitByKeys(values, []string{refine.opts.InputKey})

	initialInputs, err := refine.constructInitialInputs(docs[0], rest)
	if err != nil {
		return nil, err
	}

	res, err := refine.llmChain.Predict(ctx, initialInputs)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(docs); i++ {
		refineInputs, err := refine.constructRefineInputs(docs[i], res, rest)
		if err != nil {
			return nil, err
		}

		res, err = refine.refineLLMChain.Predict(ctx, refineInputs)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		refine.opts.OutputKey: strings.TrimSpace(res),
	}, nil
}

func (refine *RefineDocumentsChain) constructInitialInputs(doc schema.Document, rest map[string]any) (map[string]any, error) {
	docInfo := make(map[string]any)

	docInfo["pageContent"] = doc.PageContent

	for key, value := range doc.Metadata {
		docInfo[key] = value
	}

	inputs := util.CopyMap(rest)

	docText, err := refine.opts.DocumentPrompt.Format(docInfo)
	if err != nil {
		return nil, err
	}

	inputs[refine.opts.DocumentVariableName] = docText

	return inputs, nil
}

func (refine *RefineDocumentsChain) constructRefineInputs(doc schema.Document, lastResponse string, rest map[string]any) (map[string]any, error) {
	docInfo := make(map[string]any)

	docInfo["pageContent"] = doc.PageContent

	for key, value := range doc.Metadata {
		docInfo[key] = value
	}

	inputs := util.CopyMap(rest)

	docText, err := refine.opts.DocumentPrompt.Format(docInfo)
	if err != nil {
		return nil, err
	}

	inputs[refine.opts.DocumentVariableName] = docText
	inputs[refine.opts.InitialResponseName] = lastResponse

	return inputs, nil
}
