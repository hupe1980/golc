package llm

import (
	"context"
	"errors"
	"io"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the LLM interface.
var _ schema.LLM = (*OpenAI)(nil)

// OpenAIClient represents the interface for interacting with the OpenAI API.
type OpenAIClient interface {
	// CreateCompletionStream creates a streaming completion request with the provided completion request.
	// It returns a completion stream for receiving streamed completion responses from the OpenAI API.
	// The `CompletionStream` should be closed after use.
	CreateCompletionStream(ctx context.Context, request openai.CompletionRequest) (stream *openai.CompletionStream, err error)

	// CreateCompletion sends a completion request to the OpenAI API and returns the completion response.
	// It blocks until the response is received from the API.
	CreateCompletion(ctx context.Context, request openai.CompletionRequest) (response openai.CompletionResponse, err error)
}

// OpenAIOptions contains options for configuring the OpenAI LLM model.
type OpenAIOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`
	// ModelName is the name of the OpenAI language model to use.
	ModelName string `map:"model_name"`
	// Temperature is the sampling temperature to use during text generation.
	Temperatur float32 `map:"temperatur"`
	// MaxTokens is the maximum number of tokens to generate in the completion.
	MaxTokens int `map:"max_tokens"`
	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p"`
	// PresencePenalty penalizes repeated tokens.
	PresencePenalty float32 `map:"presence_penalty"`
	// FrequencyPenalty penalizes repeated tokens according to frequency.
	FrequencyPenalty float32 `map:"frequency_penalty"`
	// N is the number of completions to generate for each prompt.
	N int `map:"n"`
	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream"`
}

// OpenAI is an implementation of the LLM interface for the OpenAI language model.
type OpenAI struct {
	schema.Tokenizer
	client OpenAIClient
	opts   OpenAIOptions
}

// NewOpenAI creates a new OpenAI instance with the provided API key and options.
func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	client := openai.NewClient(apiKey)
	return NewOpenAIFromClient(client, optFns...)
}

// NewOpenAIFromClient creates a new OpenAI instance with the provided client and options.
func NewOpenAIFromClient(client OpenAIClient, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:        "text-davinci-002",
		Temperatur:       0.7,
		MaxTokens:        256,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		N:                1,
		Stream:           false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.Tokenizer == nil {
		opts.Tokenizer = tokenizer.NewOpenAI(opts.ModelName)
	}

	return &OpenAI{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *OpenAI) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	choices := []openai.CompletionChoice{}
	tokenUsage := make(map[string]int)

	completionRequest := openai.CompletionRequest{
		Prompt:           prompt,
		Model:            l.opts.ModelName,
		Temperature:      l.opts.Temperatur,
		MaxTokens:        l.opts.MaxTokens,
		TopP:             l.opts.TopP,
		PresencePenalty:  l.opts.PresencePenalty,
		FrequencyPenalty: l.opts.FrequencyPenalty,
		N:                l.opts.N,
		Stop:             opts.Stop,
	}

	if l.opts.Stream {
		completionRequest.Stream = true

		stream, err := l.client.CreateCompletionStream(ctx, completionRequest)
		if err != nil {
			return nil, err
		}

		defer stream.Close()

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

				if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
					Token: res.Choices[0].Text,
				}); err != nil {
					return nil, err
				}

				choices = append(choices, res.Choices...)
			}
		}
	} else {
		res, err := l.client.CreateCompletion(ctx, completionRequest)
		if err != nil {
			return nil, err
		}

		choices = res.Choices

		tokenUsage["CompletionTokens"] += res.Usage.CompletionTokens
		tokenUsage["PromptTokens"] += res.Usage.PromptTokens
		tokenUsage["TotalTokens"] += res.Usage.TotalTokens
	}

	generations := util.Map(choices, func(choice openai.CompletionChoice, _ int) schema.Generation {
		return schema.Generation{
			Text: choice.Text,
			Info: map[string]any{
				"FinishReason": choice.FinishReason,
				"LogProbs":     choice.LogProbs,
			},
		}
	})

	return &schema.ModelResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"ModelName":  l.opts.ModelName,
			"TokenUsage": tokenUsage,
		},
	}, nil
}

// Type returns the type of the model.
func (l *OpenAI) Type() string {
	return "llm.OpenAI"
}

// Verbose returns the verbosity setting of the model.
func (l *OpenAI) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *OpenAI) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *OpenAI) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
