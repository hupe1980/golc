package chatmodel

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/ernie"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Ernie satisfies the ChatModel interface.
var _ schema.ChatModel = (*Ernie)(nil)

// ErnieClient is the interface for the Ernie client.
type ErnieClient interface {
	CreateChatCompletion(ctx context.Context, model string, request *ernie.ChatCompletionRequest) (*ernie.ChatCompletionResponse, error)
}

// ErnieOptions is the options struct for the Ernie chat model.
type ErnieOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// ModelName is the name of the Ernie chat model to use.
	ModelName string `map:"model_name,omitempty"`

	// Temperature is the sampling temperature to use during text generation.
	Temperature float64 `map:"temperature,omitempty"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float64 `map:"top_p,omitempty"`

	// PenaltyScore is a parameter used during text generation to apply a penalty for generating longer responses.
	PenaltyScore float64 `map:"penalty_score"`
}

// Ernie is a struct representing the Ernie language model.
type Ernie struct {
	schema.Tokenizer
	client ErnieClient
	opts   ErnieOptions
}

// NewErnie creates a new instance of the Ernie chat model.
func NewErnie(clientID, clientSecret string, optFns ...func(o *ErnieOptions)) (*Ernie, error) {
	client := ernie.New(clientID, clientSecret)

	return NewErnieFromClient(client, optFns...)
}

// NewErnieFromClient creates a new instance of the Ernie chat model from a custom ErnieClient.
func NewErnieFromClient(client ErnieClient, optFns ...func(o *ErnieOptions)) (*Ernie, error) {
	opts := ErnieOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:    "ernie-bot-turbo",
		Temperature:  0.95,
		TopP:         0.7,
		PenaltyScore: 1,
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

	return &Ernie{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided chat messages and options.
func (cm *Ernie) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	ernieMessages := make([]ernie.Message, len(messages))

	for i, message := range messages {
		switch message.Type() {
		case schema.ChatMessageTypeAI:
			ernieMessages[i] = ernie.Message{
				Role:    "assistant",
				Content: message.Content(),
			}
		case schema.ChatMessageTypeHuman:
			ernieMessages[i] = ernie.Message{
				Role:    "user",
				Content: message.Content(),
			}
		case schema.ChatMessageTypeGeneric:
			m, _ := message.(schema.GenericChatMessage)

			ernieMessages[i] = ernie.Message{
				Role:    m.Role(),
				Content: m.Content(),
			}
		default:
			return nil, fmt.Errorf("unsupported message type: %s", message.Type())
		}
	}

	res, err := cm.client.CreateChatCompletion(ctx, cm.opts.ModelName, &ernie.ChatCompletionRequest{
		Messages:     ernieMessages,
		Temperature:  cm.opts.Temperature,
		TopP:         cm.opts.TopP,
		PenaltyScore: cm.opts.PenaltyScore,
	})
	if err != nil {
		return nil, err
	}

	if res.ErrorCode != 0 {
		return nil, fmt.Errorf("ernie api error: %d", res.ErrorCode)
	}

	generation := schema.Generation{
		Text:    res.Result,
		Message: schema.NewAIChatMessage(res.Result),
	}

	tokenUsage := map[string]int{
		"PromptTokens":     res.Usage.PromptTokens,
		"CompletionTokens": res.Usage.CompletionTokens,
		"TotalTokens":      res.Usage.TotalTokens,
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{generation},
		LLMOutput: map[string]any{
			"TokenUsage": tokenUsage,
		},
	}, nil
}

// Type returns the type of the model.
func (cm *Ernie) Type() string {
	return "chatmodel.Ernie"
}

// Verbose returns the verbosity setting of the model.
func (cm *Ernie) Verbose() bool {
	return cm.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Ernie) Callbacks() []schema.Callback {
	return cm.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Ernie) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}
