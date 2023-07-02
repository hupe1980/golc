package chain

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/memory"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/outputparser"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

const conversationTemplate = `The following is a friendly conversation between a human and an AI. The AI is talkative and provides lots of specific details from its context. If the AI does not know the answer to a question, it truthfully says it does not know.

Current conversation:
{{.history}}
Human: {{.input}}
AI:`

// Compile time check to ensure Conversation satisfies the Chain interface.
var _ schema.Chain = (*Conversation)(nil)

type ConversationOptions struct {
	*schema.CallbackOptions
	Prompt       *prompt.Template
	Memory       schema.Memory
	OutputKey    string
	OutputParser schema.OutputParser[any]
	// ReturnFinalOnly determines whether to return only the final parsed result or include extra generation information.
	// When set to true (default), the field will return only the final parsed result.
	// If set to false, the field will include additional information about the generation along with the final parsed result.
	ReturnFinalOnly bool
}

type Conversation struct {
	llm  schema.LLM
	opts ConversationOptions
}

func NewConversation(llm schema.LLM, optFns ...func(o *ConversationOptions)) (*Conversation, error) {
	opts := ConversationOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		OutputKey:       "response",
		Memory:          memory.NewConversationBuffer(),
		ReturnFinalOnly: true,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.OutputParser == nil {
		opts.OutputParser = outputparser.NewNoOpt()
	}

	if opts.Prompt == nil {
		opts.Prompt = prompt.NewTemplate(conversationTemplate)
	}

	return &Conversation{
		llm:  llm,
		opts: opts,
	}, nil
}

// Call executes the ConversationalRetrieval chain with the given context and inputs.
// It returns the outputs of the chain or an error, if any.
func (c *Conversation) Call(ctx context.Context, inputs schema.ChainValues, optFns ...func(o *schema.CallOptions)) (schema.ChainValues, error) {
	opts := schema.CallOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	promptValue, err := c.opts.Prompt.FormatPrompt(inputs)
	if err != nil {
		return nil, err
	}

	if opts.CallbackManger != nil {
		text := fmt.Sprintf("Prompt after formatting:\n%s", promptValue.String())
		if cbErr := opts.CallbackManger.OnText(ctx, &schema.TextManagerInput{
			Text: text,
		}); cbErr != nil {
			return nil, cbErr
		}
	}

	res, err := model.GeneratePrompt(ctx, c.llm, []schema.PromptValue{promptValue}, func(o *model.Options) {
		o.Stop = opts.Stop

		if opts.CallbackManger != nil {
			o.Callbacks = opts.CallbackManger.GetInheritableCallbacks()
			o.ParentRunID = opts.CallbackManger.RunID()
		}
	})
	if err != nil {
		return nil, err
	}

	outputs, err := c.createOutputs(res)
	if err != nil {
		return nil, err
	}

	return outputs[0], nil
}

func (c *Conversation) Prompt() *prompt.Template {
	return c.opts.Prompt
}

// Memory returns the memory associated with the chain.
func (c *Conversation) Memory() schema.Memory {
	return c.opts.Memory
}

// Type returns the type of the chain.
func (c *Conversation) Type() string {
	return "Conversation"
}

// Verbose returns the verbosity setting of the chain.
func (c *Conversation) Verbose() bool {
	return c.opts.CallbackOptions.Verbose
}

// Callbacks returns the callbacks associated with the chain.
func (c *Conversation) Callbacks() []schema.Callback {
	return c.opts.CallbackOptions.Callbacks
}

// InputKeys returns the expected input keys.
func (c *Conversation) InputKeys() []string {
	return []string{"input"}
}

// OutputKeys returns the output keys the chain will return.
func (c *Conversation) OutputKeys() []string {
	return []string{c.opts.OutputKey}
}

func (c *Conversation) createOutputs(llmResult *schema.ModelResult) ([]map[string]any, error) {
	result := make([]map[string]any, len(llmResult.Generations)-1)

	for _, generation := range llmResult.Generations {
		parsed, err := c.opts.OutputParser.ParseResult(generation)
		if err != nil {
			return nil, err
		}

		output := map[string]any{
			c.opts.OutputKey: parsed,
			"fullGeneration": generation,
		}

		result = append(result, output)
	}

	if c.opts.ReturnFinalOnly {
		for i := range result {
			result[i] = map[string]any{
				c.opts.OutputKey: result[i][c.opts.OutputKey],
			}
		}
	}

	return result, nil
}
