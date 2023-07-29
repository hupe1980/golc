package chatmodel

import (
	"context"

	generativelanguagepb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
)

// PalmClient is the interface for the PALM client.
type PalmClient interface {
	GenerateMessage(ctx context.Context, req *generativelanguagepb.GenerateMessageRequest, opts ...gax.CallOption) (*generativelanguagepb.GenerateMessageResponse, error)
}

// PalmOptions is the options struct for the PALM chat model.
type PalmOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// ModelName is the name of the Palm chat model to use.
	ModelName string `map:"model_name,omitempty"`

	// Temperature is the sampling temperature to use during text generation.
	Temperature float32 `map:"temperature,omitempty"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p,omitempty"`

	// TopK determines how the model selects tokens for output.
	TopK int32 `map:"top_k"`

	// CandidateCount specifies the number of candidates to generate during text completion.
	CandidateCount int32 `map:"candidate_count"`
}

// Palm is a struct representing the PALM language model.
type Palm struct {
	client PalmClient
	opts   PalmOptions
}

// NewPalm creates a new instance of the PALM language model.
func NewPalm(client PalmClient, optFns ...func(o *PalmOptions)) (*Palm, error) {
	opts := PalmOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:      "models/chat-bison-001",
		Temperature:    0.7,
		CandidateCount: 1,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewGPT2()
		if tErr != nil {
			return nil, tErr
		}
	}

	return &Palm{
		client: client,
		opts:   opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (cm *Palm) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	res, err := cm.client.GenerateMessage(ctx, &generativelanguagepb.GenerateMessageRequest{
		Prompt:         &generativelanguagepb.MessagePrompt{},
		Model:          cm.opts.ModelName,
		Temperature:    &cm.opts.Temperature,
		TopP:           &cm.opts.TopP,
		TopK:           &cm.opts.TopK,
		CandidateCount: &cm.opts.CandidateCount,
	})
	if err != nil {
		return nil, err
	}

	generations := util.Map(res.GetCandidates(), func(m *generativelanguagepb.Message, _ int) schema.Generation {
		switch m.GetAuthor() {
		case "ai":
			return schema.Generation{
				Message: schema.NewAIChatMessage(m.GetContent()),
				Text:    m.GetContent(),
			}
		case "human":
			return schema.Generation{
				Message: schema.NewHumanChatMessage(m.GetContent()),
				Text:    m.GetContent(),
			}
		default:
			return schema.Generation{
				Message: schema.NewGenericChatMessage(m.GetContent(), m.GetAuthor()),
				Text:    m.GetContent(),
			}
		}
	})

	return &schema.ModelResult{
		Generations: generations,
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (cm *Palm) Type() string {
	return "llm.Palm"
}

// Verbose returns the verbosity setting of the model.
func (cm *Palm) Verbose() bool {
	return cm.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Palm) Callbacks() []schema.Callback {
	return cm.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Palm) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}
