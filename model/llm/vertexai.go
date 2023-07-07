package llm

import (
	"context"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tokenizer"
	"github.com/hupe1980/golc/util"
	"google.golang.org/protobuf/types/known/structpb"
)

// Compile time check to ensure VertexAI satisfies the LLM interface.
var _ schema.LLM = (*VertexAI)(nil)

// VertexAIClient represents the interface for interacting with Vertex AI.
type VertexAIClient interface {
	// Predict sends a prediction request to the Vertex AI service.
	// It takes a context, predict request, and optional call options.
	// It returns the predict response or an error if the prediction fails.
	Predict(ctx context.Context, req *aiplatformpb.PredictRequest, opts ...gax.CallOption) (*aiplatformpb.PredictResponse, error)
}

// VertexAIOptions contains options for configuring the VertexAI language model.
type VertexAIOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Temperature is the sampling temperature to use during text generation.
	Temperatur float32 `map:"temperatur"`

	// MaxOutputTokens determines the maximum amount of text output from one prompt.
	MaxOutputTokens int `map:"max_output_tokens"`

	// TopP is the total probability mass of tokens to consider at each step.
	TopP float32 `map:"top_p"`

	// TopK determines how the model selects tokens for output.
	TopK int `map:"top_k"`
}

// VertexAI represents the VertexAI language model.
type VertexAI struct {
	schema.Tokenizer
	client   VertexAIClient
	endpoint string
	opts     VertexAIOptions
}

// NewVertexAI creates a new VertexAI instance with the provided client and endpoint.
func NewVertexAI(client VertexAIClient, endpoint string, optFns ...func(o *VertexAIOptions)) (*VertexAI, error) {
	opts := VertexAIOptions{
		Temperatur:      0.0,
		MaxOutputTokens: 128,
		TopP:            0.95,
		TopK:            40,
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

	return &VertexAI{
		Tokenizer: opts.Tokenizer,
		client:    client,
		endpoint:  endpoint,
		opts:      opts,
	}, nil
}

// Generate generates text based on the provided prompt and options.
func (l *VertexAI) Generate(ctx context.Context, prompt string, optFns ...func(o *schema.GenerateOptions)) (*schema.ModelResult, error) {
	instance, err := structpb.NewValue(map[string]any{
		"content": prompt,
	})
	if err != nil {
		return nil, err
	}

	parameters, err := structpb.NewValue(map[string]any{
		"temperature":       l.opts.Temperatur,
		"max_output_tokens": l.opts.MaxOutputTokens,
		"top_p":             l.opts.TopP,
		"top_k":             l.opts.TopK,
	})
	if err != nil {
		return nil, err
	}

	res, err := l.client.Predict(ctx, &aiplatformpb.PredictRequest{
		Endpoint:   l.endpoint,
		Instances:  []*structpb.Value{instance},
		Parameters: parameters,
	})
	if err != nil {
		return nil, err
	}

	generations := util.Map(res.Predictions, func(p *structpb.Value, _ int) schema.Generation {
		value := p.GetStructValue().AsMap()
		text, _ := value["content"].(string)

		return schema.Generation{
			Text: text,
		}
	})

	return &schema.ModelResult{
		Generations: generations,
		LLMOutput: map[string]any{
			"DeployedModelID": res.DeployedModelId,
			"Model":           res.Model,
			"ModelVersionID":  res.ModelVersionId,
			"ModelName":       res.ModelDisplayName,
		},
	}, nil
}

// Type returns the type of the model.
func (l *VertexAI) Type() string {
	return "llm.VertexAI"
}

// Verbose returns the verbosity setting of the model.
func (l *VertexAI) Verbose() bool {
	return l.opts.Verbose
}

// Callbacks returns the registered callbacks of the model.
func (l *VertexAI) Callbacks() []schema.Callback {
	return l.opts.Callbacks
}

// InvocationParams returns the parameters used in the model invocation.
func (l *VertexAI) InvocationParams() map[string]any {
	return nil
}
