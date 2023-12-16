package llm

import (
	"context"
	"errors"

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

// Compile time check to ensure Cohere satisfies the LLM interface.
var _ schema.LLM = (*Cohere)(nil)

// CohereClient is an interface for the Cohere client.
type CohereClient interface {
	Generate(ctx context.Context, request *cohere.GenerateRequest) (*cohere.Generation, error)
}

// CohereOptions contains options for configuring the Cohere LLM model.
type CohereOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model represents the name or identifier of the Cohere language model to use.
	Model string `map:"model,omitempty"`

	// NumGenerations denotes the maximum number of generations that will be returned.                   string
	NumGenerations int `map:"num_generations"`

	// MaxTokens denotes the number of tokens to predict per generation.
	MaxTokens int `map:"max_tokens"`

	// Temperature is a non-negative float that tunes the degree of randomness in generation.
	Temperature float64 `map:"temperature"`

	// K specifies the number of top most likely tokens to consider for generation at each step.
	K int `map:"k"`

	// P is a probability value between 0.0 and 1.0. It ensures that only the most likely tokens,
	// with a total probability mass of P, are considered for generation at each step.
	P float64 `map:"p"`

	// FrequencyPenalty is used to reduce repetitiveness of generated tokens. A higher value applies
	// a stronger penalty to previously present tokens, proportional to how many times they have
	// already appeared in the prompt or prior generation.
	FrequencyPenalty float64 `map:"frequency_penalty"`

	// PresencePenalty is used to reduce repetitiveness of generated tokens. It applies a penalty
	// equally to all tokens that have already appeared, regardless of their exact frequencies.
	PresencePenalty float64 `map:"presence_penalty"`

	// ReturnLikelihoods specifies whether and how the token likelihoods are returned with the response.
	// It can be set to "GENERATION", "ALL", or "NONE". If "GENERATION" is selected, the token likelihoods
	// will only be provided for generated text. If "ALL" is selected, the token likelihoods will be
	// provided for both the prompt and the generated text.
	ReturnLikelihoods string `map:"return_likelihoods,omitempty"`

	// MaxRetries represents the maximum number of retries to make when generating.
	MaxRetries uint `map:"max_retries,omitempty"`
}

// Cohere represents the Cohere language model.
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
		Model:             "medium",
		NumGenerations:    1,
		MaxTokens:         256,
		Temperature:       0.75,
		K:                 0,
		P:                 1,
		FrequencyPenalty:  0,
		PresencePenalty:   0,
		ReturnLikelihoods: "NONE",
		MaxRetries:        3,
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

// Generate generates text based on the provided prompt and options.
func (l *Cohere) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	returnLikelihoods, err := cohere.NewGenerateRequestReturnLikelihoodsFromString(l.opts.ReturnLikelihoods)
	if err != nil {
		return nil, err
	}

	res, err := l.generateWithRetry(ctx, &cohere.GenerateRequest{
		Model:             util.AddrOrNil(l.opts.Model),
		NumGenerations:    util.AddrOrNil(l.opts.NumGenerations),
		MaxTokens:         util.AddrOrNil(l.opts.MaxTokens),
		Temperature:       util.AddrOrNil(l.opts.Temperature),
		K:                 util.AddrOrNil(l.opts.K),
		P:                 util.AddrOrNil(l.opts.P),
		PresencePenalty:   util.AddrOrNil(l.opts.PresencePenalty),
		FrequencyPenalty:  util.AddrOrNil(l.opts.FrequencyPenalty),
		ReturnLikelihoods: returnLikelihoods.Ptr(),
		Prompt:            prompt,
		StopSequences:     opts.Stop,
	})
	if err != nil {
		return nil, err
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: res.Generations[0].Text}},
		LLMOutput: map[string]any{
			"likelihood":       res.Generations[0].Likelihood,
			"tokenLikelihoods": res.Generations[0].TokenLikelihoods,
		},
	}, nil
}

func (l *Cohere) generateWithRetry(ctx context.Context, req *cohere.GenerateRequest) (*cohere.Generation, error) {
	retryOpts := []retry.Option{
		retry.Attempts(l.opts.MaxRetries),
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

	var res *cohere.Generation

	err := retry.Do(
		func() error {
			r, cErr := l.client.Generate(ctx, req)
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
func (l *Cohere) Type() string {
	return "llm.Cohere"
}

// Verbose returns the verbosity setting of the model.
func (l *Cohere) Verbose() bool {
	return l.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *Cohere) Callbacks() []schema.Callback {
	return l.opts.CallbackOptions.Callbacks
}

// InvocationParams returns the parameters used in the llm model invocation.
func (l *Cohere) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
