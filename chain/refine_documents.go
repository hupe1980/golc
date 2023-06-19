package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

type RefineDocumentsOptions struct {
	*callbackOptions
	InputKey             string
	DocumentVariableName string
	InitialResponseName  string
	DocumentPrompt       *prompt.Template
	OutputKey            string
}

type RefineDocumentsChain struct {
	*baseChain
	llmChain       *LLMChain
	refineLLMChain *LLMChain
	opts           RefineDocumentsOptions
}

func NewRefineDocumentsChain(llmChain *LLMChain, refineLLMChain *LLMChain, optFns ...func(o *RefineDocumentsOptions)) (*RefineDocumentsChain, error) {
	opts := RefineDocumentsOptions{
		InputKey:             "inputDocuments",
		DocumentVariableName: "context",
		InitialResponseName:  "existingAnswer",
		OutputKey:            "text",
		callbackOptions: &callbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
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

	refine.baseChain = &baseChain{
		chainName:       "RefineDocumentsChain",
		callFunc:        refine.call,
		inputKeys:       []string{opts.InputKey},
		outputKeys:      llmChain.OutputKeys(),
		callbackOptions: opts.callbackOptions,
	}

	return refine, nil
}

func (c *RefineDocumentsChain) call(ctx context.Context, values schema.ChainValues) (schema.ChainValues, error) {
	input, ok := values[c.opts.InputKey]
	if !ok {
		return nil, fmt.Errorf("%w: no value for inputKey %s", ErrInvalidInputValues, c.opts.InputKey)
	}

	docs, ok := input.([]schema.Document)
	if !ok {
		return nil, ErrInputValuesWrongType
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("%w: documents slice has no elements", ErrInvalidInputValues)
	}

	rest := util.OmitByKeys(values, []string{c.opts.InputKey})

	initialInputs, err := c.constructInitialInputs(docs[0], rest)
	if err != nil {
		return nil, err
	}

	res, err := c.llmChain.Predict(ctx, initialInputs)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(docs); i++ {
		refineInputs, err := c.constructRefineInputs(docs[i], res, rest)
		if err != nil {
			return nil, err
		}

		res, err = c.refineLLMChain.Predict(ctx, refineInputs)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		c.opts.OutputKey: strings.TrimSpace(res),
	}, nil
}

func (c *RefineDocumentsChain) Memory() schema.Memory {
	return nil
}

func (c *RefineDocumentsChain) Type() string {
	return "RefineDocumentsChain"
}

func (c *RefineDocumentsChain) Verbose() bool {
	return c.opts.callbackOptions.Verbose
}

func (c *RefineDocumentsChain) Callbacks() []schema.Callback {
	return c.opts.callbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *RefineDocumentsChain) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *RefineDocumentsChain) OutputKeys() []string {
	return c.llmChain.OutputKeys()
}

func (c *RefineDocumentsChain) constructInitialInputs(doc schema.Document, rest map[string]any) (map[string]any, error) {
	docInfo := make(map[string]any)

	docInfo["pageContent"] = doc.PageContent
	for key, value := range doc.Metadata {
		docInfo[key] = value
	}

	inputs := util.CopyMap(rest)

	docText, err := c.opts.DocumentPrompt.Format(docInfo)
	if err != nil {
		return nil, err
	}

	inputs[c.opts.DocumentVariableName] = docText

	return inputs, nil
}

func (c *RefineDocumentsChain) constructRefineInputs(doc schema.Document, lastResponse string, rest map[string]any) (map[string]any, error) {
	docInfo := make(map[string]any)

	docInfo["pageContent"] = doc.PageContent

	for key, value := range doc.Metadata {
		docInfo[key] = value
	}

	inputs := util.CopyMap(rest)

	docText, err := c.opts.DocumentPrompt.Format(docInfo)
	if err != nil {
		return nil, err
	}

	inputs[c.opts.DocumentVariableName] = docText
	inputs[c.opts.InitialResponseName] = lastResponse

	return inputs, nil
}
