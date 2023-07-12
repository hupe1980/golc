package embedding

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hupe1980/go-tiktoken"
	"github.com/hupe1980/golc/util"
	"github.com/sashabaranov/go-openai"
)

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
}

type OpenAI struct {
	client *openai.Client
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

func NewOpenAIFromClient(client *openai.Client, optFns ...func(o *OpenAIOptions)) (*OpenAI, error) {
	opts := OpenAIOptions{
		ModelName:              "text-embedding-ada-002",
		EmbeddingContextLength: 8191,
		ChunkSize:              1000,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &OpenAI{
		client: client,
		opts:   opts,
	}, nil
}

// EmbedDocuments embeds a list of documents and returns their embeddings.
func (e *OpenAI) EmbedDocuments(ctx context.Context, texts []string) ([][]float64, error) {
	return e.getLenSafeEmbeddings(ctx, texts)
}

// EmbedQuery embeds a single query and returns its embedding.
func (e *OpenAI) EmbedQuery(ctx context.Context, text string) ([]float64, error) {
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

	res, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: nameToOpenAIModel[e.opts.ModelName],
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}

	return util.Map(res.Data[0].Embedding, func(e float32, i int) float64 {
		return float64(e)
	}), nil
}

func (e *OpenAI) getLenSafeEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
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

	batchedEmbeddings := [][]float64{}

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
			batchedEmbeddings = append(batchedEmbeddings, util.Map(d.Embedding, func(e float32, _ int) float64 {
				return float64(e)
			}))
		}
	}

	results := make([][][]float64, len(texts))
	numTokensInBatch := make([][]int, len(texts))

	for i := 0; i < len(indices); i++ {
		index := indices[i]
		results[index] = append(results[index], batchedEmbeddings[i])
		numTokensInBatch[index] = append(numTokensInBatch[index], len(tokens[i]))
	}

	embeddings := make([][]float64, len(texts))

	for i := 0; i < len(texts); i++ {
		var average []float64

		result := results[i]

		if len(result) == 0 {
			res, err := e.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
				Model: nameToOpenAIModel[e.opts.ModelName],
				Input: []string{""},
			})
			if err != nil {
				return nil, err
			}

			average = util.Map(res.Data[0].Embedding, func(e float32, i int) float64 {
				return float64(e)
			})
		} else {
			sum := make([]float64, len(result[0]))

			weights := numTokensInBatch[i]

			for j := 0; j < len(result); j++ {
				embedding := result[j]
				for k := 0; k < len(embedding); k++ {
					sum[k] += embedding[k] * float64(weights[j])
				}
			}

			average = make([]float64, len(sum))
			for j := 0; j < len(sum); j++ {
				average[j] = sum[j] / float64(util.SumInt(weights))
			}
		}

		norm := 0.0
		for _, value := range average {
			norm += value * value
		}

		norm = math.Sqrt(norm)
		for j := 0; j < len(average); j++ {
			average[j] /= norm
		}

		embeddings[i] = average
	}

	return embeddings, nil
}
