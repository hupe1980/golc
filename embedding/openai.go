package embedding

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hupe1980/go-tiktoken"
	"github.com/hupe1980/golc/internal/math32"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/sashabaranov/go-openai"
)

// Compile time check to ensure OpenAI satisfies the Embedder interface.
var _ schema.Embedder = (*OpenAI)(nil)

type OpenAIClient interface {
	CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequestConverter) (res openai.EmbeddingResponse, err error)
}

// nolint staticcheck
var nameToOpenAIModel = map[string]openai.EmbeddingModel{
	"text-similarity-ada-001":       openai.AdaSimilarity,
	"text-similarity-babbage-001":   openai.BabbageSimilarity,
	"text-similarity-curie-001":     openai.CurieSimilarity,
	"text-similarity-davinci-001":   openai.DavinciSimilarity,
	"text-search-ada-doc-001":       openai.AdaSearchDocument,
	"text-search-ada-query-001":     openai.AdaSearchQuery,
	"text-search-babbage-doc-001":   openai.BabbageSearchDocument,
	"text-search-babbage-query-001": openai.BabbageSearchQuery,
	"text-search-curie-doc-001":     openai.CurieSearchDocument,
	"text-search-curie-query-001":   openai.CurieSearchQuery,
	"text-search-davinci-doc-001":   openai.DavinciSearchDocument,
	"text-search-davinci-query-001": openai.DavinciSearchQuery,
	"code-search-ada-code-001":      openai.AdaCodeSearchCode,
	"code-search-ada-text-001":      openai.AdaCodeSearchText,
	"code-search-babbage-code-001":  openai.BabbageCodeSearchCode,
	"code-search-babbage-text-001":  openai.BabbageCodeSearchText,
	"text-embedding-ada-002":        openai.AdaEmbeddingV2,
}

type OpenAIOptions struct {
	// Model name to use.
	ModelName              string
	EmbeddingContextLength int
	// Maximum number of texts to embed in each batch
	ChunkSize int
	// BaseURL is the base URL of the OpenAI service.
	BaseURL string
	// OrgID is the organization ID for accessing the OpenAI service.
	OrgID string
	// MaxRetries represents the maximum number of retries to make when embedding.
	MaxRetries uint `map:"max_retries,omitempty"`
}

var DefaultOpenAIConfig = OpenAIOptions{
	ModelName:              "text-embedding-ada-002",
	EmbeddingContextLength: 8191,
	ChunkSize:              1000,
	MaxRetries:             3,
}

type OpenAI struct {
	client OpenAIClient
	opts   OpenAIOptions
}

func NewOpenAI(apiKey string, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{}

	for _, fn := range optFns {
		fn(&opts)
	}

	config := openai.DefaultConfig(apiKey)

	if opts.BaseURL != "" {
		config.BaseURL = opts.BaseURL
	}

	if opts.OrgID != "" {
		config.OrgID = opts.OrgID
	}

	client := openai.NewClientWithConfig(config)

	return NewOpenAIFromClient(client, optFns...)
}

func NewOpenAIFromClient(client OpenAIClient, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := DefaultOpenAIConfig

	for _, fn := range optFns {
		fn(&opts)
	}

	return &OpenAI{
		client: client,
		opts:   opts,
	}, nil
}

// BatchEmbedText embeds a list of texts and returns their embeddings.
func (e *OpenAI) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	return e.getLenSafeEmbeddings(ctx, texts)
}

// EmbedText embeds a single text and returns its embedding.
func (e *OpenAI) EmbedText(ctx context.Context, text string) ([]float32, error) {
	if len(text) > e.opts.EmbeddingContextLength {
		embeddings, err := e.getLenSafeEmbeddings(ctx, []string{text})
		if err != nil {
			return nil, err
		}

		return embeddings[0], nil
	}

	if strings.HasSuffix(e.opts.ModelName, "001") {
		// See: https://github.com/openai/openai-python/issues/418#issuecomment-1525939500
		// replace newlines, which can negatively affect performance.
		text = strings.ReplaceAll(text, "\n", " ")
	}

	res, err := e.createEmbeddingsWithRetry(ctx, openai.EmbeddingRequest{
		Model: nameToOpenAIModel[e.opts.ModelName],
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}

	return res.Data[0].Embedding, nil
}

func (e *OpenAI) createEmbeddingsWithRetry(ctx context.Context, request openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error) {
	retryOpts := []retry.Option{
		retry.Attempts(e.opts.MaxRetries),
		retry.DelayType(retry.FixedDelay),
		retry.RetryIf(func(err error) bool {
			e := &openai.APIError{}
			if errors.As(err, &e) {
				switch e.HTTPStatusCode {
				case 429, 500:
					return true
				default:
					return false
				}
			}

			return false
		}),
	}

	var res openai.EmbeddingResponse

	err := retry.Do(
		func() error {
			r, cErr := e.client.CreateEmbeddings(ctx, request)
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

func (e *OpenAI) getLenSafeEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// please refer to
	// https://github.com/openai/openai-cookbook/blob/main/examples/Embedding_long_inputs.ipynb
	tokens := []string{}
	indices := []int{}

	encoding, err := tiktoken.NewEncodingForModel(e.opts.ModelName)
	if err != nil {
		return nil, err
	}

	for i, text := range texts {
		if strings.HasSuffix(e.opts.ModelName, "001") {
			// Replace newlines, which can negatively affect performance.
			text = strings.ReplaceAll(text, "\n", " ")
		}

		token, _, err := encoding.Encode(text, nil, nil)
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(token); j += e.opts.EmbeddingContextLength {
			limit := j + e.opts.EmbeddingContextLength
			if limit > len(token) {
				limit = len(token)
			}

			tokens = append(tokens, util.Map(token[j:limit], func(e uint, _ int) string {
				return fmt.Sprintf("%d", e)
			})...)

			indices = append(indices, i)
		}
	}

	batchedEmbeddings := [][]float32{}

	for i := 0; i < len(tokens); i += e.opts.ChunkSize {
		limit := i + e.opts.ChunkSize
		if limit > len(tokens) {
			limit = len(tokens)
		}

		res, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Model: nameToOpenAIModel[e.opts.ModelName],
			Input: tokens[i:limit],
		})
		if err != nil {
			return nil, err
		}

		for _, d := range res.Data {
			batchedEmbeddings = append(batchedEmbeddings, d.Embedding)
		}
	}

	results := make([][][]float32, len(texts))
	numTokensInBatch := make([][]int, len(texts))

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		results[index] = append(results[index], batchedEmbeddings[i])
		numTokensInBatch[index] = append(numTokensInBatch[index], len(tokens[i]))
	}

	embeddings := make([][]float32, len(texts))

	for i := 0; i < len(texts); i++ {
		var average []float32

		result := results[i]

		if len(result) == 0 {
			res, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
				Model: nameToOpenAIModel[e.opts.ModelName],
				Input: []string{""},
			})
			if err != nil {
				return nil, err
			}

			average = res.Data[0].Embedding
		} else {
			sum := make([]float32, len(result[0]))

			weights := numTokensInBatch[i]

			for j := 0; j < len(result); j++ {
				embedding := result[j]
				for k := 0; k < len(embedding); k++ {
					sum[k] += embedding[k] * float32(weights[j])
				}
			}

			average = make([]float32, len(sum))
			for j := 0; j < len(sum); j++ {
				average[j] = sum[j] / float32(util.SumInt(weights))
			}
		}

		norm := float32(0.0)
		for _, value := range average {
			norm += value * value
		}

		norm = math32.Sqrt(norm)
		for j := 0; j < len(average); j++ {
			average[j] /= norm
		}

		embeddings[i] = average
	}

	return embeddings, nil
}
