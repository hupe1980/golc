package documentcompressor

import (
	"context"
	"testing"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestCohereRerank(t *testing.T) {
	t.Parallel()

	t.Run("Documents without metadata", func(t *testing.T) {
		t.Parallel()

		// Arrange
		mockClient := &mockCohereClient{
			rerankResponse: &cohere.RerankResponse{
				Results: []*cohere.RerankResponseResultsItem{
					{Index: 0, RelevanceScore: 0.8},
					{Index: 1, RelevanceScore: 0.6},
				},
			},
		}

		compressor := NewCohereRank(mockClient)

		// Input documents
		docs := []schema.Document{
			{PageContent: "Document 1"},
			{PageContent: "Document 2"},
		}

		// Test
		result, err := compressor.Compress(context.Background(), docs, "query")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 0.8, result[0].Metadata["relevanceScore"])
		assert.Equal(t, 0.6, result[1].Metadata["relevanceScore"])
	})

	t.Run("Documents with metadata", func(t *testing.T) {
		t.Parallel()

		// Arrange
		mockClient := &mockCohereClient{
			rerankResponse: &cohere.RerankResponse{
				Results: []*cohere.RerankResponseResultsItem{
					{Index: 0, RelevanceScore: 0.8},
					{Index: 1, RelevanceScore: 0.6},
				},
			},
		}

		compressor := NewCohereRank(mockClient)

		// Input documents
		docs := []schema.Document{
			{PageContent: "Document 1", Metadata: map[string]any{"foo": "bar"}},
			{PageContent: "Document 2", Metadata: map[string]any{"foo": "bar"}},
		}

		// Test
		result, err := compressor.Compress(context.Background(), docs, "query")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, 0.8, result[0].Metadata["relevanceScore"])
		assert.Equal(t, "bar", result[0].Metadata["foo"])
		assert.Equal(t, 0.6, result[1].Metadata["relevanceScore"])
		assert.Equal(t, "bar", result[1].Metadata["foo"])
	})
}

// mockCohereClient is a custom mock implementation of the CohereClient interface.
type mockCohereClient struct {
	rerankResponse *cohere.RerankResponse
	rerankErr      error
}

// Rerank is a mock implementation of the Rerank method.
func (m *mockCohereClient) Rerank(ctx context.Context, request *cohere.RerankRequest) (*cohere.RerankResponse, error) {
	if m.rerankErr != nil {
		return nil, m.rerankErr
	}

	return m.rerankResponse, nil
}
