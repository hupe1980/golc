package chatmodel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/avast/retry-go"
	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	core "github.com/cohere-ai/cohere-go/v2/core"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Cohere satisfies the ChatModel interface.
var _ schema.ChatModel = (*Cohere)(nil)

// CohereClient defines the interface for interacting with the Cohere API.
type CohereClient interface {
	// Chat performs a non-streaming chat with the Cohere API, generating a response based on the provided request.
	// It returns a non-streamed response or an error if the operation fails.
	Chat(ctx context.Context, request *cohere.ChatRequest, opts ...core.RequestOption) (*cohere.NonStreamedChatResponse, error)

	// ChatStream performs a streaming chat with the Cohere API, generating responses in a stream based on the provided request.
	// It returns a stream of responses or an error if the operation fails.
	ChatStream(ctx context.Context, request *cohere.ChatStreamRequest, opts ...core.RequestOption) (*core.Stream[cohere.StreamedChatResponse], error)
}

// CohereOptions contains options for configuring the Cohere model.
type CohereOptions struct {
	// CallbackOptions specify options for handling callbacks during text generation.
	*schema.CallbackOptions `map:"-"`

	// Tokenizer represents the tokenizer to be used with the LLM model.
	schema.Tokenizer `map:"-"`

	// Model represents the name or identifier of the Cohere language model to use.
	Model string `map:"model,omitempty"`

	// Temperature is a non-negative float that tunes the degree of randomness in generation.
	Temperature float64 `map:"temperature"`

	// MaxRetries represents the maximum number of retries to make when generating.
	MaxRetries uint `map:"max_retries,omitempty"`

	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// Cohere represents an instance of the Cohere language model.
type Cohere struct {
	schema.Tokenizer
	client CohereClient
	opts   CohereOptions
}

// NewCohere creates a new Cohere instance using the provided API key and optional configuration options.
// It internally creates a Cohere client using the provided API key and initializes the Cohere struct.
func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	client := cohereclient.NewClient(cohereclient.WithToken(apiKey))
	return NewCohereFromClient(client, optFns...)
}

// NewCohereFromClient creates a new Cohere instance using the provided Cohere client and optional configuration options.
func NewCohereFromClient(client CohereClient, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		Model:       "command",
		Temperature: 0.75,
		MaxRetries:  3,
		Stream:      false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		var tErr error

		opts.Tokenizer, tErr = tokenizer.NewCohere("coheretext-50k")
		if tErr != nil {
			return nil, tErr
		}
	}

	return &Cohere{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided chat messages and options.
func (cm *Cohere) Generate(ctx context.Context, messages schema.ChatMessages, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("at least one message must be passed")
	}

	chatMessages := make([]*cohere.ChatMessage, len(messages)-1)

	for i, m := range messages[:len(messages)-1] {
		var role cohere.ChatMessageRole

		switch m.Type() {
		case schema.ChatMessageTypeAI:
			role = cohere.ChatMessageRoleChatbot
		case schema.ChatMessageTypeHuman:
			role = cohere.ChatMessageRoleUser
		default:
			return nil, fmt.Errorf("unsupported chat message type: %s", m.Type())
		}

		chatMessages[i] = &cohere.ChatMessage{
			Role:    role,
			Message: m.Content(),
		}
	}

	var text string

	if cm.opts.Stream {
		stream, err := cm.client.ChatStream(ctx, &cohere.ChatStreamRequest{
			Model:       util.AddrOrNil(cm.opts.Model),
			Message:     messages[len(messages)-1].Content(),
			ChatHistory: chatMessages,
			Temperature: util.AddrOrNil(cm.opts.Temperature),
		})
		if err != nil {
			return nil, err
		}

		defer stream.Close()

		var tokens []string

	streamProcessing:
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				res, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break streamProcessing
				}
				if err != nil {
					return nil, err
				}

				if res.EventType == "text-generation" {
					if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
						Token: res.TextGeneration.Text,
					}); err != nil {
						return nil, err
					}

					tokens = append(tokens, res.TextGeneration.Text)
				}
			}
		}

		text = strings.Join(tokens, "")
	} else {
		res, err := cm.generateWithRetry(ctx, &cohere.ChatRequest{
			Model:       util.AddrOrNil(cm.opts.Model),
			Message:     messages[len(messages)-1].Content(),
			ChatHistory: chatMessages,
			Temperature: util.AddrOrNil(cm.opts.Temperature),
		})
		if err != nil {
			return nil, err
		}

		text = res.Text
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{newChatGeneraton(text)},
		LLMOutput:   map[string]any{},
	}, nil
}

func (cm *Cohere) generateWithRetry(ctx context.Context, req *cohere.ChatRequest) (*cohere.NonStreamedChatResponse, error) {
	retryOpts := []retry.Option{
		retry.Attempts(cm.opts.MaxRetries),
		retry.DelayType(retry.FixedDelay),
		retry.RetryIf(func(err error) bool {
			e := new(core.APIError)
			if errors.As(err, &e) {
				switch e.StatusCode {
				case 429, 500:
					return true
				default:
					return false
				}
			}

			return false
		}),
	}

	var res *cohere.NonStreamedChatResponse

	err := retry.Do(
		func() error {
			r, cErr := cm.client.Chat(ctx, req)
			if cErr != nil {
				return cErr
			}

			res = r

			return nil
		},
		retryOpts...,
	)

	return res, err
}

// Type returns the type of the model.
func (cm *Cohere) Type() string {
	return "chatmodel.Cohere"
}

// Verbose returns the verbosity setting of the model.
func (cm *Cohere) Verbose() bool {
	return cm.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (cm *Cohere) Callbacks() []schema.Callback {
	return cm.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (cm *Cohere) InvocationParams() map[string]any {
	return util.StructToMap(cm.opts)
}
