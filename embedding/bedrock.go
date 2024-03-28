package embedding

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/schema"
	"golang.org/x/sync/errgroup"
)

// Compile time check to ensure Bedrock satisfies the Embedder interface.
var _ schema.Embedder = (*Bedrock)(nil)

// BedrockInputOutputAdapter is a helper struct for preparing input and handling output for Bedrock model.
type BedrockInputOutputAdapter struct {
	provider string
}

// NewBedrockInputOutputAdpter creates a new instance of BedrockInputOutputAdpter.
func NewBedrockInputOutputAdapter(provider string) *BedrockInputOutputAdapter {
	return &BedrockInputOutputAdapter{
		provider: provider,
	}
}

// PrepareInput prepares the input for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareInput(text string, modelParams map[string]any) ([]byte, error) {
	var body map[string]any

	text = removeNewLines(text)

	switch bioa.provider {
	case "amazon":
		body = make(map[string]any)
		body["inputText"] = text
	case "cohere":
		body = modelParams

		if _, ok := body["input_type"]; !ok {
			body["input_type"] = "search_document"
		}

		body["texts"] = []string{text}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", bioa.provider)
	}

	return json.Marshal(body)
}

// amazonOutput represents the expected JSON output structure from the Bedrock model for the Amazon provider..
type amazonOutput struct {
	Embedding []float32 `json:"embedding"`
}

// cohereOutput represents the expected JSON output structure from the Bedrock model for the Cohere provider.
type cohereOutput struct {
	Embeddings []float32 `json:"embeddings"`
}

// PrepareOutput prepares the output for the Bedrock model based on the specified provider.
func (bioa *BedrockInputOutputAdapter) PrepareOutput(response []byte) ([]float32, error) {
	switch bioa.provider {
	case "amazon":
		output := &amazonOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return nil, err
		}

		return output.Embedding, nil
	case "cohere":
		output := &cohereOutput{}
		if err := json.Unmarshal(response, output); err != nil {
			return nil, err
		}

		return output.Embeddings, nil
	}

	return nil, fmt.Errorf("unsupported provider: %s", bioa.provider)
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// BedrockAmazonOptions is a struct containing options for configuring the Amazon Bedrock model.
type BedrockAmazonOptions struct {
	// Model id to use.
	ModelID string `map:"model_id,omitempty"`
}

// NewBedrockAmazon creates a new instance of Bedrock with the Amazon provider.
func NewBedrockAmazon(client BedrockRuntimeClient, optFns ...func(o *BedrockAmazonOptions)) *Bedrock {
	opts := BedrockAmazonOptions{
		ModelID: "amazon.titan-embed-text-v1",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return NewBedrock(client, opts.ModelID)
}

// BedrockCohereOptions is a struct containing options for configuring the Cohere Bedrock model.
type BedrockCohereOptions struct {
	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	InputType string `map:"input_type"`

	Truncate string `map:"truncate"`
}

// NewBedrockCohere creates a new instance of Bedrock with the Cohere provider.
func NewBedrockCohere(client BedrockRuntimeClient, optFns ...func(o *BedrockCohereOptions)) *Bedrock {
	opts := BedrockCohereOptions{
		ModelID:   "cohere.embed-english-v3",
		InputType: "search_document",
		Truncate:  "NONE",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return NewBedrock(client, opts.ModelID, func(o *BedrockOptions) {
		o.ModelParams = map[string]interface{}{
			"input_type": opts.InputType,
			"truncate":   opts.Truncate,
		}
	})
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	MaxConcurrency int

	// Model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`
}

// Bedrock is a struct representing the Bedrock model embedding functionality.
type Bedrock struct {
	client BedrockRuntimeClient

	// Model id to use.
	modelID string `map:"model_id,omitempty"`

	opts BedrockOptions
}

// NewBedrock creates a new instance of Bedrock with the provided BedrockRuntimeClient and optional configuration.
func NewBedrock(client BedrockRuntimeClient, modelID string, optFns ...func(o *BedrockOptions)) *Bedrock {
	opts := BedrockOptions{
		MaxConcurrency: 5,
		ModelParams:    make(map[string]any),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Bedrock{
		client:  client,
		modelID: modelID,
		opts:    opts,
	}
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Bedrock) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	errs, errctx := errgroup.WithContext(ctx)

	errs.SetLimit(e.opts.MaxConcurrency)

	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		i, text := i, text

		errs.Go(func() error {
			embedding, err := e.EmbedText(errctx, text)
			if err != nil {
				return err
			}

			embeddings[i] = embedding

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return embeddings, nil
}

// EmbedText embeds a single text and returns its embedding.
func (e *Bedrock) EmbedText(ctx context.Context, text string) ([]float32, error) {
	bioa := NewBedrockInputOutputAdapter(e.getProvider())

	body, err := bioa.PrepareInput(text, e.opts.ModelParams)
	if err != nil {
		return nil, err
	}

	res, err := e.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(e.modelID),
		Body:        body,
		Accept:      aws.String("application/json"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	return bioa.PrepareOutput(res.Body)
}

// getProvider returns the provider of the model based on the model ID.
func (e *Bedrock) getProvider() string {
	return strings.Split(e.modelID, ".")[0]
}
