package rag

import (
	"context"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure RefineDocuments satisfies the Chain interface.
var _ schema.Chain = (*RefineDocuments)(nil)

type RefineDocumentsOptions struct {
	*schema.CallbackOptions
	InputKey             string
	DocumentVariableName string
	InitialResponseName  string
	DocumentPrompt       schema.PromptTemplate
	OutputKey            string
}

type RefineDocuments struct {
	llmChain       *chain.LLM
	refineLLMChain *chain.LLM
	opts           RefineDocumentsOptions
}

func NewRefineDocuments(llmChain *chain.LLM, refineLLMChain *chain.LLM, optFns ...func(o *RefineDocumentsOptions)) (*RefineDocuments, error) {
	opts := RefineDocumentsOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		InputKey:             "inputDocuments",
		OutputKey:            "outputText",
		DocumentVariableName: "text",
		InitialResponseName:  "existingAnswer",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.DocumentPrompt == nil {
		opts.DocumentPrompt = prompt.NewTemplate("{{.pageContent}}")
	}

	return &RefineDocuments{
		llmChain:       llmChain,
		refineLLMChain: refineLLMChain,
		opts:           opts,
	}, nil
}

// Call executes the RefineDocuments chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *RefineDocuments) Call(ctx context.Context, values schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	docs, err := values.GetDocuments(c.opts.InputKey)
	if err != nil {
		return nil, err
	}

	rest := util.OmitByKeys(values, []string{c.opts.InputKey})

	initialInputs, err := c.constructInitialInputs(docs[0], rest)
	if err != nil {
		return nil, err
	}

	res, err := golc.SimpleCall(ctx, c.llmChain, initialInputs, func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
	})
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(docs); i++ {
		refineInputs, err := c.constructRefineInputs(docs[i], res, rest)
		if err != nil {
			return nil, err
		}

		res, err = golc.SimpleCall(ctx, c.refineLLMChain, refineInputs, func(sco *golc.SimpleCallOptions) {
			sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
			sco.ParentRunID = opts.CallbackManger.RunID()
		})
		if err != nil {
			return nil, err
		}
	}

	return schema.ChainValues{
		c.opts.OutputKey: strings.TrimSpace(res),
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *RefineDocuments) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *RefineDocuments) Type() string {
	return "RefineDocuments"
}

// Verbose returns the verbosity setting of the chain.
func (c *RefineDocuments) Verbose() bool {
	return c.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *RefineDocuments) Callbacks() []schema.Callback {
	return c.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (c *RefineDocuments) InputKeys() []string {
	return []string{c.opts.InputKey}
}

// OutputKeys returns the output keys the chain will return.
func (c *RefineDocuments) OutputKeys() []string {
	return c.llmChain.OutputKeys()
}

func (c *RefineDocuments) constructInitialInputs(doc schema.Document, rest map[string]any) (schema.ChainValues, error) {
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

func (c *RefineDocuments) constructRefineInputs(doc schema.Document, lastResponse string, rest map[string]any) (schema.ChainValues, error) {
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
