package embedding

import (
	"context"
	"errors"

	"github.com/avast/retry-go"
	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	core "github.com/cohere-ai/cohere-go/v2/core"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Cohere satisfies the Embedder interface.
var _ schema.Embedder = (*Cohere)(nil)

// CohereClient is an interface for the Cohere client.
type CohereClient interface {
	Embed(ctx context.Context, request *cohere.EmbedRequest) (*cohere.EmbedResponse, error)
}

// CohereOptions contains options for configuring the Cohere instance.
type CohereOptions struct {
	// Model name to use.
	Model string
	// Truncate embeddings that are too long from start or end ("NONE"|"START"|"END")
	Truncate string
	// MaxRetries represents the maximum number of retries to make when embedding.
	MaxRetries uint `map:"max_retries,omitempty"`
}

// Cohere is a client for the Cohere API.
type Cohere struct {
	client CohereClient
	opts   CohereOptions
}

// NewCohere creates a new Cohere instance with the provided API key and options.
// It returns the initialized Cohere instance or an error if initialization fails.
func NewCohere(apiKey string, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	client := cohereclient.NewClient(cohereclient.WithToken(apiKey))

	return NewCohereFromClient(client, optFns...)
}

// NewCohereFromClient creates a new Cohere instance from an existing Cohere client and options.
// It returns the initialized Cohere instance.
func NewCohereFromClient(client CohereClient, optFns ...func(o *CohereOptions)) (*Cohere, error) {
	opts := CohereOptions{
		Model:      "embed-english-v3.0",
		MaxRetries: 3,
		Truncate:   "NONE",
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Cohere{
		client: client,
		opts:   opts,
	}, nil
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *Cohere) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	res, err := e.embedWithRetry(ctx, &cohere.EmbedRequest{
		Model:    util.AddrOrNil(e.opts.Model),
		Truncate: cohere.EmbedRequestTruncate(e.opts.Truncate).Ptr(),
		Texts:    texts,
	})
	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(res.Embeddings))
	for i, r := range res.Embeddings {
		embeddings[i] = util.Float64ToFloat32(r)
	}

	return embeddings, nil
}

func (e *Cohere) embedWithRetry(ctx context.Context, req *cohere.EmbedRequest) (*cohere.EmbedResponse, error) {
	retryOpts := []retry.Option{
		retry.Attempts(e.opts.MaxRetries),
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

	var res *cohere.EmbedResponse

	err := retry.Do(
		func() error {
			r, cErr := e.client.Embed(ctx, req)
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

// EmbedText embeds a single query and returns its embedding.
func (e *Cohere) EmbedText(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := e.BatchEmbedText(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return embeddings[0], nil
}
