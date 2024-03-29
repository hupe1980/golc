package llm

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure Ollama satisfies the LLM interface.
var _ schema.LLM = (*Ollama)(nil)

// OllamaClient is an interface for the Ollama generative model client.
type OllamaClient interface {
	// CreateGeneration produces a single request and response for the Ollama generative model.
	CreateGeneration(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationResponse, error)
	// CreateGenerationStream initiates a streaming request and returns a stream for the Ollama generative model.
	CreateGenerationStream(ctx context.Context, req *ollama.GenerationRequest) (*ollama.GenerationStream, error)
}

// OllamaOptions contains options for the Ollama model.
type OllamaOptions struct {
	// CallbackOptions specify options for handling callbacks during text generation.
	*schema.CallbackOptions `map:"-"`
	// Tokenizer represents the tokenizer to be used with the LLM model.
	schema.Tokenizer `map:"-"`
	// ModelName is the name of the Gemini model to use.
	ModelName string `map:"model_name,omitempty"`
	// Temperature controls the randomness of the generation. Higher values make the output more random.
	Temperature float32 `map:"temperature,omitempty"`
	// MaxTokens is the maximum number of tokens to generate in the completion.
	MaxTokens int `map:"max_tokens,omitempty"`
	// TopP is the nucleus sampling parameter. It controls the cumulative probability of the most likely tokens to sample from.
	TopP float32 `map:"top_p,omitempty"`
	// TopK is the number of top tokens to consider for sampling.
	TopK int `map:"top_k,omitempty"`
	// PresencePenalty penalizes repeated tokens.
	PresencePenalty float32 `map:"presence_penalty,omitempty"`
	// FrequencyPenalty penalizes repeated tokens according to frequency.
	FrequencyPenalty float32 `map:"frequency_penalty,omitempty"`
	// Stream indicates whether to stream the results or not.
	Stream bool `map:"stream,omitempty"`
}

// Ollama is a struct representing the Ollama generative model.
type Ollama struct {
	schema.Tokenizer
	client OllamaClient
	opts   OllamaOptions
}

// NewOllama creates a new instance of the Ollama model with the provided client and options.
func NewOllama(client OllamaClient, optFns ...func(o *OllamaOptions)) (*Ollama, error) {
	opts := OllamaOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:        "llama2",
		Temperature:      0.7,
		MaxTokens:        256,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
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

	return &Ollama{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *Ollama) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	req := &ollama.GenerationRequest{
		Model:  l.opts.ModelName,
		Prompt: prompt,
		Options: ollama.Options{
			Temperature:      l.opts.Temperature,
			NumPredict:       l.opts.MaxTokens,
			TopK:             l.opts.TopK,
			TopP:             l.opts.TopP,
			PresencePenalty:  l.opts.PresencePenalty,
			FrequencyPenalty: l.opts.FrequencyPenalty,
			Stop:             opts.Stop,
		},
	}

	var text string

	if l.opts.Stream {
		req.Stream = util.PTR(true)

		stream, err := l.client.CreateGenerationStream(ctx, req)
		if err != nil {
			return nil, err
		}

		defer stream.Close()

		tokens := []string{}

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

				if !res.Done {
					if err := opts.CallbackManger.OnModelNewToken(ctx, &schema.ModelNewTokenManagerInput{
						Token: res.Response,
					}); err != nil {
						return nil, err
					}

					tokens = append(tokens, res.Response)
				}
				// else {
				// 	// TODO Metrics, EvalCount, ... -> LLMOutput?
				// }
			}

			text = strings.Join(tokens, "")
		}
	} else {
		res, err := l.client.CreateGeneration(ctx, req)
		if err != nil {
			return nil, err
		}

		text = res.Response
	}

	return &schema.ModelResult{
		Generations: []schema.Generation{{Text: text}},
		LLMOutput:   map[string]any{
			//"Done": res.Done,
		},
	}, nil
}

// Type returns the type of the model.
func (l *Ollama) Type() string {
	return "llm.Ollama"
}

// Verbose returns the verbosity setting of the model.
func (l *Ollama) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *Ollama) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *Ollama) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
