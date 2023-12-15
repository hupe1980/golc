package embedding

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/schema"
	"golang.org/x/sync/errgroup"
)

// Compile time check to ensure Bedrock satisfies the Embedder interface.
var _ schema.Embedder = (*Bedrock)(nil)

// amazonOutput represents the expected JSON output structure from the Bedrock model.
type amazonOutput struct {
	Embedding []float32 `json:"embedding"`
}

// BedrockRuntimeClient is an interface for the Bedrock model runtime client.
type BedrockRuntimeClient interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
}

// BedrockOptions contains options for configuring the Bedrock model.
type BedrockOptions struct {
	*schema.CallbackOptions `map:"-"`
	schema.Tokenizer        `map:"-"`

	MaxConcurrency int

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
		MaxConcurrency: 5,
		ModelID:        "amazon.titan-embed-text-v1",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Bedrock{
		client: client,
		opts:   opts,
	}
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Bedrock) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	errs, errctx := errgroup.WithContext(ctx)

	// Use a semaphore to control concurrency
	sem := make(chan struct{}, e.opts.MaxConcurrency)

	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		// Acquire semaphore, limit concurrency
		sem <- struct{}{}

		i, text := i, text

		errs.Go(func() error {
			defer func() { <-sem }() // Release semaphore when done

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
	jsonBody := map[string]string{
		"inputText": removeNewLines(text),
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
