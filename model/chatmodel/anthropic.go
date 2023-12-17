package chatmodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/anthropic"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

const (
	humanPromptPrefix = "\n\nHuman:"
	aiPromptPrefix    = "\n\nAssistant:"
)

// Compile time check to ensure Anthropic satisfies the ChatModel interface.
var _ schema.ChatModel = (*Anthropic)(nil)

// AnthropicClient is the interface for the Anthropic client.
type AnthropicClient interface {
	CreateCompletion(ctx context.Context, request *anthropic.CompletionRequest) (*anthropic.CompletionResponse, error)
}

// AnthropicOptions contains options for configuring the Anthropic chat model.
type AnthropicOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model name to use.
	ModelName string `map:"model_name,omitempty"`

	// Temperature parameter controls the randomness of the generation output.
	Temperature float32 `map:"temperature,omitempty"`

	// Denotes the number of tokens to predict per generation.
	MaxTokens int `map:"max_tokens,omitempty"`

	// TopK parameter specifies the number of highest probability tokens to consider for generation.
	TopK int `map:"top_k,omitempty"`

	// TopP parameter specifies the cumulative probability threshold for generating tokens.
	TopP float32 `map:"top_p,omitempty"`
}

// Anthropic is a chat model based on the Anthropic API.
type Anthropic struct {
	schema.Tokenizer
	client AnthropicClient
	opts   AnthropicOptions
}

// NewAnthropic creates a new instance of the Anthropic chat model with the provided options.
func NewAnthropic(apiKey string, optFns ...func(o *AnthropicOptions)) (*Anthropic, error) {
	client := anthropic.New(apiKey)

	return NewAnthropicFromClient(client, optFns...)
}

// NewAnthropicFromClient creates a new instance of the Anthropic chat model with the provided options.
func NewAnthropicFromClient(client AnthropicClient, optFns ...func(o *AnthropicOptions)) (*Anthropic, error) {
	opts := AnthropicOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:   "claude-v1",
		Temperature: 0.5,
		MaxTokens:   256,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewClaude()
		if tErr != nil {
			return nil, tErr
		}
	}

	return &Anthropic{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided chat messages and options.
func (cm *Anthropic) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	prompt, err := convertMessagesToAnthropicPrompt(messages)
	if err != nil {
		return nil, err
	}

	res, err := cm.client.CreateCompletion(ctx, &anthropic.CompletionRequest{
		Prompt:      prompt,
		Model:       cm.opts.ModelName,
		Temperature: cm.opts.Temperature,
		MaxTokens:   cm.opts.MaxTokens,
		TopK:        cm.opts.TopK,
		TopP:        cm.opts.TopP,
		Stop:        opts.Stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{newChatGeneraton(res.Completion)},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (cm *Anthropic) Type() string {
	return "chatmodel.Anthropic"
}

// Verbose returns the verbosity setting of the model.
func (cm *Anthropic) Verbose() bool {
	return cm.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Anthropic) Callbacks() []schema.Callback {
	return cm.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Anthropic) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}

func convertMessagesToAnthropicPrompt(messages schema.ChatMessages) (string, error) {
	if len(messages) > 0 {
		msg := messages[len(messages)-1]
		if msg.Type() != schema.ChatMessageTypeAI {
			messages = append(messages, schema.NewAIChatMessage(""))
		}
	}

	prompt := ""

	for _, message := range messages {
		switch message.Type() {
		case schema.ChatMessageTypeSystem:
			prompt += fmt.Sprintf("%s <admin>%s</admin>", humanPromptPrefix, message.Content())
		case schema.ChatMessageTypeAI:
			prompt += fmt.Sprintf("%s %s", aiPromptPrefix, message.Content())
		case schema.ChatMessageTypeHuman:
			prompt += fmt.Sprintf("%s %s", humanPromptPrefix, message.Content())
		default:
			return "", fmt.Errorf("unsupported message type: %s", message.Type())
		}
	}

	return strings.TrimRight(prompt, " "), nil
}
