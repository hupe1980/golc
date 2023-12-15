package llm

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/ai/generativelanguage/apiv1/generativelanguagepb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
)

// Compile time check to ensure GoogleGenAI satisfies the LLM interface.
var _ schema.LLM = (*GoogleGenAI)(nil)

// GoogleGenAIClient is an interface for the GoogleGenAI model client.
type GoogleGenAIClient interface {
	GenerateContent(context.Context, *generativelanguagepb.GenerateContentRequest, ...gax.CallOption) (*generativelanguagepb.GenerateContentResponse, error)
	CountTokens(context.Context, *generativelanguagepb.CountTokensRequest, ...gax.CallOption) (*generativelanguagepb.CountTokensResponse, error)
}

// GoogleGenAIOptions contains options for the GoogleGenAI Language Model.
type GoogleGenAIOptions struct {
	// CallbackOptions specify options for handling callbacks during text generation.
	*schema.CallbackOptions `map:"-"`
	// Tokenizer represents the tokenizer to be used with the LLM model.
	schema.Tokenizer `map:"-"`
	// ModelName is the name of the GoogleGenAI model to use.
	ModelName string `map:"model_name,omitempty"`
	// CandidateCount is the number of candidate generations to consider.
	CandidateCount int32 `map:"candidate_count,omitempty"`
	// MaxOutputTokens is the maximum number of tokens to generate in the output.
	MaxOutputTokens int32 `map:"max_output_tokens,omitempty"`
	// Temperature controls the randomness of the generation. Higher values make the output more random.
	Temperature float32 `map:"temperature,omitempty"`
	// TopP is the nucleus sampling parameter. It controls the cumulative probability of the most likely tokens to sample from.
	TopP float32 `map:"top_p,omitempty"`
	// TopK is the number of top tokens to consider for sampling.
	TopK int32 `map:"top_k,omitempty"`
}

// GoogleGenAI represents the GoogleGenAI Language Model.
type GoogleGenAI struct {
	schema.Tokenizer
	client GoogleGenAIClient
	opts   GoogleGenAIOptions
}

// NewGoogleGenAI creates a new instance of the GoogleGenAI Language Model.
func NewGoogleGenAI(client GoogleGenAIClient, optFns ...func(o *GoogleGenAIOptions)) (*GoogleGenAI, error) {
	opts := GoogleGenAIOptions{
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
		ModelName:       "models/gemini-pro",
		CandidateCount:  1,
		MaxOutputTokens: 2048,
		TopK:            3,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if !strings.HasPrefix(opts.ModelName, "models/") {
		opts.ModelName = fmt.Sprintf("models/%s", opts.ModelName)
	}

	if opts.Tokenizer == nil {
		opts.Tokenizer = tokenizer.NewGoogleGenAI(client, opts.ModelName)
	}

	return &GoogleGenAI{
		Tokenizer: opts.Tokenizer,
		client:    client,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *GoogleGenAI) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	opts := schema.GenerateOptions{
		CallbackManger: &callback.NoopManager{},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	res, err := l.client.GenerateContent(ctx, &generativelanguagepb.GenerateContentRequest{
		Model: l.opts.ModelName,
		Contents: []*generativelanguagepb.Content{{Parts: []*generativelanguagepb.Part{{
			Data: &generativelanguagepb.Part_Text{Text: prompt},
		}}}},
		GenerationConfig: &generativelanguagepb.GenerationConfig{
			CandidateCount:  util.AddrOrNil(l.opts.CandidateCount),
			MaxOutputTokens: util.AddrOrNil(l.opts.MaxOutputTokens),
			Temperature:     util.AddrOrNil(l.opts.Temperature),
			TopP:            util.AddrOrNil(l.opts.TopP),
			TopK:            util.AddrOrNil(l.opts.TopK),
			StopSequences:   opts.Stop,
		},
	})
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
func (l *GoogleGenAI) Type() string {
	return "llm.GoogleGenAI"
}

// Verbose returns the verbosity setting of the model.
func (l *GoogleGenAI) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *GoogleGenAI) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *GoogleGenAI) InvocationParams() map[string]any {
	return util.StructToMap(l.opts)
}
