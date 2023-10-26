package embedding

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Bedrock satisfies the Embedder interface.
var _ schema.Embedder = (*Bedrock)(nil)

// amazonOutput represents the expected JSON output structure from the Bedrock model.
type amazonOutput struct {
	Embedding []float64 `json:"embedding"`
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	// Model id to use.
	ModelID string `map:"model_id,omitempty"`

	// Model params to use.
	ModelParams map[string]any `map:"model_params,omitempty"`
}

// Bedrock is a struct representing the Bedrock model embedding functionality.
type Bedrock struct {
	client BedrockRuntimeClient
	opts   BedrockOptions
}

// NewBedrock creates a new instance of Bedrock with the provided BedrockRuntimeClient and optional configuration.
func NewBedrock(client BedrockRuntimeClient, optFns ...func(o *BedrockOptions)) *Bedrock {
	opts := BedrockOptions{
		ModelID: "amazon.titan-embed-text-v1",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Bedrock{
		client: client,
		opts:   opts,
	}
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *Bedrock) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))

	for i, text := range texts {
		embedding, err := e.EmbedQuery(ctx, text)
		if err != nil {
			return nil, err
		}

		embeddings[i] = embedding
	}

	return embeddings, nil
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *Bedrock) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
	jsonBody := map[string]string{
		"inputText": text,
	}

	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	res, err := e.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(e.opts.ModelID),
		Body:        body,
		Accept:      aws.String("application/json"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, err
	}

	output := &amazonOutput{}
	if err := json.Unmarshal(res.Body, output); err != nil {
		return nil, err
	}

	return output.Embedding, nil
}
