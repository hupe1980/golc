package embedding

import (
	"context"
	"testing"

	huggingface "github.com/hupe1980/go-huggingface"
	"github.com/stretchr/testify/assert"
)

func TestHuggingFaceHub(t *testing.T) {
	// Create a mock client with responses.
	mockResponses := map[string][]float64{
		"document1": {0.1, 0.2, 0.3},
		"document2": {0.4, 0.5, 0.6},
		"query1":    {0.7, 0.8, 0.9},
	}
	mockClient := &mockHuggingFaceHubClient{
		Responses: mockResponses,
	}

	// Create an instance of HuggingFaceHub.
	embedder := NewHuggingFaceHubFromClient(mockClient)

	t.Run("EmbedDocuments", func(t *testing.T) {
		// Define test documents.
		documents := []string{"document1", "document2"}

		// Expected embeddings for the test documents.
		expectedEmbeddings := [][]float64{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
		}

		// Test the EmbedDocuments method.
		embeddings, err := embedder.EmbedDocuments(context.Background(), documents)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmbeddings, embeddings)
	})

	t.Run("EmbedQuery", func(t *testing.T) {
		// Define a test query.
		query := "query1"

		// Expected embedding for the test query.
		expectedEmbedding := []float64{0.7, 0.8, 0.9}

		// Test the EmbedQuery method.
		embedding, err := embedder.EmbedQuery(context.Background(), query)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmbedding, embedding)
	})
}

// mockHuggingFaceHubClient is a mock implementation of the HuggingFaceHubClient interface for testing.
type mockHuggingFaceHubClient struct {
	Responses map[string][]float64
	Err       error
}

func (m *mockHuggingFaceHubClient) FeatureExtractionWithAutomaticReduction(ctx context.Context, req *huggingface.FeatureExtractionRequest) (huggingface.FeatureExtractionWithAutomaticReductionResponse, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	var resp huggingface.FeatureExtractionWithAutomaticReductionResponse

	// Simulate model responses based on the provided inputs (query text).
	for _, i := range req.Inputs {
		if embeddings, ok := m.Responses[i]; ok {
			resp = append(resp, embeddings)
		}
	}

	return resp, nil
}
