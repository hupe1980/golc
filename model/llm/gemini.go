package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"

	"github.com/google/generative-ai-go/genai"
)

// Compile time check to ensure Gemini satisfies the LLM interface.
var _ schema.LLM = (*Gemini)(nil)

// GeminiClient is an interface for the Gemini model client.
type GeminiClient interface {
	GenerativeModel(name string) *genai.GenerativeModel
}

// GeminiOptions contains options for the Gemini Language Model.
type GeminiOptions struct {
	// CallbackOptions specify options for handling callbacks during text generation.
	*schema.CallbackOptions `map:"-"`
	// Tokenizer represents the tokenizer to be used with the LLM model.
	schema.Tokenizer `map:"-"`
	// ModelName is the name of the Gemini model to use.
	ModelName string `map:"model_name,omitempty"`
	// CandidateCount is the number of candidate generations to consider.
	CandidateCount int32
	// MaxOutputTokens is the maximum number of tokens to generate in the output.
	MaxOutputTokens int32
	// Temperature controls the randomness of the generation. Higher values make the output more random.
	Temperature float32
	// TopP is the nucleus sampling parameter. It controls the cumulative probability of the most likely tokens to sample from.
	TopP float32
	// TopK is the number of top tokens to consider for sampling.
	TopK int32
}

// Gemini represents the Gemini Language Model.
type Gemini struct {
	schema.Tokenizer
	client GeminiClient
	opts   GeminiOptions
}

// NewGemini creates a new instance of the Gemini Language Model.
func NewGemini(client GeminiClient, optFns ...func(o *GeminiOptions)) (*Gemini, error) {
	opts := GeminiOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:       "gemini-pro",
		CandidateCount:  1,
		MaxOutputTokens: 2048,
		TopK:            3,
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

	return &Gemini{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *Gemini) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	model := l.client.GenerativeModel(l.opts.ModelName)

	model.GenerationConfig = genai.GenerationConfig{
		CandidateCount:  l.opts.CandidateCount,
		MaxOutputTokens: l.opts.MaxOutputTokens,
		Temperature:     l.opts.Temperature,
		TopP:            l.opts.TopP,
		TopK:            l.opts.TopK,
		StopSequences:   opts.Stop,
	}

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	generations := make([]schema.Generation, len(res.Candidates))

	for i, c := range res.Candidates {
		var b strings.Builder
		for _, p := range c.Content.Parts {
			fmt.Fprintf(&b, "%v", p)
		}

		generations[i] = schema.Generation{Text: b.String()}
	}

	return &schema.ModelResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"BlockReason": res.PromptFeedback.BlockReason.String(),
		},
	}, nil
}

// Type returns the type of the model.
func (l *Gemini) Type() string {
	return "llm.Gemini"
}

// Verbose returns the verbosity setting of the model.
func (l *Gemini) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *Gemini) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *Gemini) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
