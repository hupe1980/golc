package llm

import (
	"context"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/ai21"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure AI21 satisfies the LLM interface.
var _ schema.LLM = (*AI21)(nil)

// AI21Client is an interface for interacting with the AI21 API.
type AI21Client interface {
	CreateCompletion(ctx context.Context, model string, req *ai21.CompleteRequest) (*ai21.CompleteResponse, error)
}

// DefaultPenalty represents the default penalty options for AI21 text completion.
var DefaultPenalty = ai21.Penalty{
	Scale:               0,
	ApplyToWhitespaces:  true,
	ApplyToPunctuations: true,
	ApplyToNumbers:      true,
	ApplyToStopwords:    true,
	ApplyToEmojis:       true,
}

// AI21Options contains options for configuring the AI21 LLM model.
type AI21Options struct {
	// CallbackOptions specify options for handling callbacks during text generation.
	*schema.CallbackOptions `map:"-"`
	// Tokenizer represents the tokenizer to be used with the LLM model.
	schema.Tokenizer `map:"-"`

	// Model is the name of the AI21 model to use for text completion.
	Model string `map:"model,omitempty"`

	// Temperature controls the randomness of text generation. Higher values make it more random.
	Temperature float64 `map:"temperature"`

	// MaxTokens sets the maximum number of tokens in the generated text.
	MaxTokens int `map:"maxTokens"`

	// MinTokens sets the minimum number of tokens in the generated text.
	MinTokens int `map:"minTokens"`

	// TopP sets the nucleus sampling probability. Higher values result in more diverse text.
	TopP float64 `map:"topP"`

	// PresencePenalty specifies the penalty for repeating words in generated text.
	PresencePenalty ai21.Penalty `map:"presencePenalty"`

	// CountPenalty specifies the penalty for repeating tokens in generated text.
	CountPenalty ai21.Penalty `map:"countPenalty"`

	// FrequencyPenalty specifies the penalty for generating frequent words.
	FrequencyPenalty ai21.Penalty `map:"frequencyPenalty"`

	// NumResults sets the number of completion results to return.
	NumResults int `map:"numResults"`
}

// AI21 is an AI21 LLM model that generates text based on a provided response function.
type AI21 struct {
	schema.Tokenizer
	client AI21Client
	opts   AI21Options
}

// NewAI21 creates a new AI21 LLM instance with the provided API key and optional configuration options.
func NewAI21(apiKey string, optFns ...func(o *AI21Options)) (*AI21, error) {
	client := ai21.New(apiKey)
	return NewAI21FromClient(client, optFns...)
}

// NewAI21FromClient creates a new AI21 LLM instance using the provided AI21 client and optional configuration options.
func NewAI21FromClient(client AI21Client, optFns ...func(o *AI21Options)) (*AI21, error) {
	opts := AI21Options{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		Model:            "j2-mid",
		Temperature:      0.7,
		MaxTokens:        256,
		MinTokens:        0,
		TopP:             1,
		PresencePenalty:  DefaultPenalty,
		CountPenalty:     DefaultPenalty,
		FrequencyPenalty: DefaultPenalty,
		NumResults:       1,
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

	return &AI21{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *AI21) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	res, err := l.client.CreateCompletion(ctx, l.opts.Model, &ai21.CompleteRequest{
		Prompt:           prompt,
		Temperature:      l.opts.Temperature,
		MaxTokens:        l.opts.MaxTokens,
		MinTokens:        l.opts.MinTokens,
		TopP:             l.opts.TopP,
		PresencePenalty:  l.opts.PresencePenalty,
		CountPenalty:     l.opts.CountPenalty,
		FrequencyPenalty: l.opts.FrequencyPenalty,
		NumResults:       l.opts.NumResults,
		StopSequences:    opts.Stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: res.Completions[0].Data.Text}},
		LLMOutput:   map[string]any{},
	}, nil
}

// Type returns the type of the model.
func (l *AI21) Type() string {
	return "llm.AI21"
}

// Verbose returns the verbosity setting of the model.
func (l *AI21) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *AI21) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *AI21) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
