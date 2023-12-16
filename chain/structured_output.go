package chain

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/jsonschema"
	"github.com/hupe1980/golc/internal/deepcopy"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StructuredOutput satisfies the Chain interface.
var _ schema.Chain = (*StructuredOutput)(nil)

const defaultStructuredOutputTemplate = `Extract and save the relevant entities mentioned in the following passage.
		
Passage:
{{.input}}`

// OutputCandidate represents a candidate for structured output containing a name,
// description, and data of any struct type.
type OutputCandidate struct {
	Name        string
	Description string
	Data        any
}

// StructuredOutputOptions contains options for configuring the StructuredOutput chain.
type StructuredOutputOptions struct {
	*schema.CallbackOptions
	Prompt    prompt.ChatTemplate
	OutputKey string
}

var DefaultStructuredOutputTemplate = StructuredOutputOptions{
	CallbackOptions: &schema.CallbackOptions{},
	Prompt: prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewHumanMessageTemplate(defaultStructuredOutputTemplate),
	}),
	OutputKey: "output",
}

// StructuredOutput is a chain that generates structured output using a ChatModel chain and candidate values.
type StructuredOutput struct {
	chatModelChain *ChatModel
	candidatesMap  map[string]OutputCandidate
	opts           StructuredOutputOptions
}

// NewStructuredOutput creates a new StructuredOutput chain with the given ChatModel, prompt, and candidates.
func NewStructuredOutput(chatModel schema.ChatModel, candidates []OutputCandidate, optFns ...func(o *StructuredOutputOptions)) (*StructuredOutput, error) {
	opts := DefaultStructuredOutputTemplate
	opts.Verbose = golc.Verbose

	for _, fn := range optFns {
		fn(&opts)
	}

	candidatesMap := make(map[string]OutputCandidate, len(candidates))
	for _, c := range candidates {
		candidatesMap[c.Name] = c
	}

	functions := make([]schema.FunctionDefinition, 0, len(candidates))

	for name, v := range candidatesMap {
		jsonSchema, err := jsonschema.Generate(reflect.TypeOf(v.Data))
		if err != nil {
			return nil, err
		}

		functions = append(functions, schema.FunctionDefinition{
			Name:        name,
			Description: v.Description,
			Parameters: schema.FunctionDefinitionParameters{
				Type:       "object",
				Properties: jsonSchema.Properties,
				Required:   jsonSchema.Required,
			},
		})
	}

	chatModelChain, err := NewChatModelWithFunctions(chatModel, opts.Prompt, functions, func(o *ChatModelOptions) {
		o.CallbackOptions = opts.CallbackOptions
		o.ForceFunctionCall = true
	})
	if err != nil {
		return nil, err
	}

	return &StructuredOutput{
		chatModelChain: chatModelChain,
		candidatesMap:  candidatesMap,
		opts:           opts,
	}, nil
}

// Call executes the StructuredOutput chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *StructuredOutput) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	output, err := golc.Call(ctx, c.chatModelChain, inputs, func(sco *golc.CallOptions) {
		sco.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
		sco.ParentRunID = opts.CallbackManger.RunID()
		sco.Stop = opts.Stop
	})
	if err != nil {
		return nil, err
	}

	aiMsg, ok := output[c.chatModelChain.OutputKeys()[0]].(*schema.AIChatMessage)
	if !ok {
		return nil, errors.New("unexpected output: message is not a ai chat message")
	}

	ext := aiMsg.Extension()
	if ext.FunctionCall == nil {
		return nil, errors.New("unexpected output: message without function call extension")
	}

	data := deepcopy.Copy(c.candidatesMap[ext.FunctionCall.Name].Data)

	if err := json.Unmarshal([]byte(ext.FunctionCall.Arguments), &data); err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: data,
	}, nil
}

// Memory returns the memory associated with the chain.
func (c *StructuredOutput) Memory() schema.Memory {
	return nil
}

// Type returns the type of the chain.
func (c *StructuredOutput) Type() string {
	return "StructuredOutput"
}

// Verbose returns the verbosity setting of the chain.
func (c *StructuredOutput) Verbose() bool {
	return c.opts.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *StructuredOutput) Callbacks() []schema.Callback {
	return c.opts.Callbacks
}

// InputKeys returns the expected input keys.
func (c *StructuredOutput) InputKeys() []string {
	return c.chatModelChain.InputKeys()
}

// OutputKeys returns the output keys the chain will return.
func (c *StructuredOutput) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}
