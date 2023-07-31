package chain

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/jsonschema"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure StructuredOutput satisfies the Chain interface.
var _ schema.Chain = (*StructuredOutput)(nil)

// OutputCandidates represents a map of candidate names to their descriptions and values used in the structured output chain.
type OutputCandidates map[string]struct {
	Description string
	Candidate   any
}

// StructuredOutputOptions contains options for configuring the StructuredOutput chain.
type StructuredOutputOptions struct {
	*schema.CallbackOptions
	OutputKey string
}

// StructuredOutput is a chain that generates structured output using a ChatModel chain and candidate values.
type StructuredOutput struct {
	chatModelChain *ChatModel
	candidates     OutputCandidates
	opts           StructuredOutputOptions
}

// NewStructuredOutput creates a new StructuredOutput chain with the given ChatModel, prompt, and candidates.
func NewStructuredOutput(chatModel schema.ChatModel, prompt prompt.ChatTemplate, candidates OutputCandidates, optFns ...func(o *StructuredOutputOptions)) (*StructuredOutput, error) {
	opts := StructuredOutputOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		OutputKey: "output",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	functions := make([]schema.FunctionDefinition, 0, len(candidates))

	for name, v := range candidates {
		jsonSchema, err := jsonschema.Generate(reflect.TypeOf(v.Candidate))
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

	chatModelChain, err := NewChatModelWithFunctions(chatModel, prompt, functions)
	if err != nil {
		return nil, err
	}

	return &StructuredOutput{
		chatModelChain: chatModelChain,
		candidates:     candidates,
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
	})
	if err != nil {
		return nil, err
	}

	aiMsg, ok := output["message"].(*schema.AIChatMessage)
	if !ok {
		return nil, errors.New("unexpected output: message is not a ai chat message")
	}

	ext := aiMsg.Extension()
	if ext.FunctionCall == nil {
		return nil, errors.New("unexpected output: message without function call extension")
	}

	out := c.candidates[ext.FunctionCall.Name]

	if err := json.Unmarshal([]byte(ext.FunctionCall.Arguments), &out.Candidate); err != nil {
		return nil, err
	}

	return schema.ChainValues{
		c.opts.OutputKey: out.Candidate,
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
