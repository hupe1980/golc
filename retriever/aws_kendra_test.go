package retriever

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

// Mocked kendra.Client for testing
type mockKendraClient struct{}

func (m *mockKendraClient) Query(ctx context.Context, params *kendra.QueryInput, optFns ...func(*kendra.Options)) (*kendra.QueryOutput, error) {
	// Simulate a successful query with some mock results
	results := []types.QueryResultItem{
		{
			DocumentTitle: &types.TextWithHighlights{
				Text: aws.String("Document 1"),
			},
			DocumentURI: aws.String("https://example.com/document1"),
			Type:        types.QueryResultTypeDocument,
			DocumentExcerpt: &types.TextWithHighlights{
				Text: aws.String("Excerpt 1"),
			},
		},
		{
			DocumentTitle: &types.TextWithHighlights{
				Text: aws.String("Document 2"),
			},
			DocumentURI: aws.String("https://example.com/document2"),
			Type:        types.QueryResultTypeDocument,
			DocumentExcerpt: &types.TextWithHighlights{
				Text: aws.String("Excerpt 2"),
			},
		},
	}

	return &kendra.QueryOutput{
		ResultItems: results,
	}, nil
}

func TestAWSKendra(t *testing.T) {
	retriever := NewAWSKendra(&mockKendraClient{}, "index", func(opts *AWSKendraOptions) {
		opts.K = 2
		opts.LanguageCode = "en"
	})

	t.Run("GetRelevantDocuments", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			query := "test query"
			expectedDocuments := []schema.Document{
				{
					PageContent: "Document Title: Document 1\nDocument Excerpt: Excerpt 1\n",
					Metadata: map[string]interface{}{
						"source":  "https://example.com/document1",
						"title":   "Document 1",
						"excerpt": "Excerpt 1",
						"type":    "DOCUMENT",
					},
				},
				{
					PageContent: "Document Title: Document 2\nDocument Excerpt: Excerpt 2\n",
					Metadata: map[string]interface{}{
						"source":  "https://example.com/document2",
						"title":   "Document 2",
						"excerpt": "Excerpt 2",
						"type":    "DOCUMENT",
					},
				},
			}

			documents, err := retriever.GetRelevantDocuments(context.Background(), query)
			assert.NoError(t, err, "GetRelevantDocuments should not return an error")
			assert.Equal(t, expectedDocuments, documents, "Retrieved documents should match expected result")
		})
	})
}
