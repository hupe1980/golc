package retriever

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestAmazonKendra_GetRelevantDocuments(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		retrieveOutput *kendra.RetrieveOutput
		retrieveError  error
		queryOutput    *kendra.QueryOutput
		queryError     error
		expectedDocs   []schema.Document
		expectedError  error
	}{
		{
			name: "Retrieve success",
			retrieveOutput: &kendra.RetrieveOutput{
				ResultItems: []types.RetrieveResultItem{
					{
						DocumentTitle: aws.String("Title 1"),
						Content:       aws.String("Content 1"),
						DocumentURI:   aws.String("URI 1"),
					},
					{
						DocumentTitle: aws.String("Title 2"),
						Content:       aws.String("Content 2"),
						DocumentURI:   aws.String("URI 2"),
					},
				},
			},
			retrieveError: nil,
			queryOutput:   nil,
			queryError:    nil,
			expectedDocs: []schema.Document{
				{
					PageContent: "Document Title: Title 1\nDocument Excerpt: Content 1\n",
					Metadata: map[string]interface{}{
						"source":  "URI 1",
						"title":   "Title 1",
						"excerpt": "Content 1",
					},
				},
				{
					PageContent: "Document Title: Title 2\nDocument Excerpt: Content 2\n",
					Metadata: map[string]interface{}{
						"source":  "URI 2",
						"title":   "Title 2",
						"excerpt": "Content 2",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Retrieve error",
			retrieveOutput: nil,
			retrieveError:  errors.New("retrieve error"),
			queryOutput:    nil,
			queryError:     nil,
			expectedDocs:   nil,
			expectedError:  errors.New("retrieve error"),
		},
		{
			name: "Retrieve no result items",
			retrieveOutput: &kendra.RetrieveOutput{
				ResultItems: []types.RetrieveResultItem{},
			},
			retrieveError: nil,
			queryOutput: &kendra.QueryOutput{
				ResultItems: []types.QueryResultItem{
					{
						DocumentTitle: &types.TextWithHighlights{
							Text: aws.String("Title 1"),
						},
						DocumentExcerpt: &types.TextWithHighlights{
							Text: aws.String("Excerpt 1"),
						},
						DocumentURI: aws.String("URI 1"),
						Type:        types.QueryResultTypeDocument,
					},
					{
						DocumentTitle: &types.TextWithHighlights{
							Text: aws.String("Title 2"),
						},
						DocumentExcerpt: &types.TextWithHighlights{
							Text: aws.String("Excerpt 2"),
						},
						DocumentURI: aws.String("URI 2"),
						Type:        types.QueryResultTypeDocument,
					},
				},
			},
			queryError: nil,
			expectedDocs: []schema.Document{
				{
					PageContent: "Document Title: Title 1\nDocument Excerpt: Excerpt 1\n",
					Metadata: map[string]interface{}{
						"source":  "URI 1",
						"title":   "Title 1",
						"excerpt": "Excerpt 1",
						"type":    "DOCUMENT",
					},
				},
				{
					PageContent: "Document Title: Title 2\nDocument Excerpt: Excerpt 2\n",
					Metadata: map[string]interface{}{
						"source":  "URI 2",
						"title":   "Title 2",
						"excerpt": "Excerpt 2",
						"type":    "DOCUMENT",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Query success",
			retrieveOutput: &kendra.RetrieveOutput{
				ResultItems: []types.RetrieveResultItem{},
			},
			retrieveError: nil,
			queryOutput: &kendra.QueryOutput{
				ResultItems: []types.QueryResultItem{
					{
						DocumentTitle: &types.TextWithHighlights{
							Text: aws.String("Title 1"),
						},
						DocumentExcerpt: &types.TextWithHighlights{
							Text: aws.String("Excerpt 1"),
						},
						DocumentURI: aws.String("URI 1"),
						Type:        types.QueryResultTypeDocument,
					},
					{
						DocumentTitle: &types.TextWithHighlights{
							Text: aws.String("Title 2"),
						},
						DocumentExcerpt: &types.TextWithHighlights{
							Text: aws.String("Excerpt 2"),
						},
						DocumentURI: aws.String("URI 2"),
						Type:        types.QueryResultTypeDocument,
					},
				},
			},
			queryError: nil,
			expectedDocs: []schema.Document{
				{
					PageContent: "Document Title: Title 1\nDocument Excerpt: Excerpt 1\n",
					Metadata: map[string]interface{}{
						"source":  "URI 1",
						"title":   "Title 1",
						"excerpt": "Excerpt 1",
						"type":    "DOCUMENT",
					},
				},
				{
					PageContent: "Document Title: Title 2\nDocument Excerpt: Excerpt 2\n",
					Metadata: map[string]interface{}{
						"source":  "URI 2",
						"title":   "Title 2",
						"excerpt": "Excerpt 2",
						"type":    "DOCUMENT",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Query error",
			retrieveOutput: &kendra.RetrieveOutput{
				ResultItems: []types.RetrieveResultItem{},
			},
			retrieveError: nil,
			queryOutput:   nil,
			queryError:    errors.New("query error"),
			expectedDocs:  nil,
			expectedError: errors.New("query error"),
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the retriever with the mock client
			r := NewAmazonKendra(&mockAmazonKendraClient{
				RetrieveOutput: tt.retrieveOutput,
				RetrieveError:  tt.retrieveError,
				QueryOutput:    tt.queryOutput,
				QueryError:     tt.queryError,
			}, "index-id")

			// Call the method under test
			docs, err := r.GetRelevantDocuments(context.Background(), "query")

			// Check the results
			assert.Equal(t, tt.expectedDocs, docs)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// mockAmazonKendraClient is a mock implementation of the AmazonKendraClient interface.
type mockAmazonKendraClient struct {
	RetrieveOutput *kendra.RetrieveOutput
	RetrieveError  error
	QueryOutput    *kendra.QueryOutput
	QueryError     error
}

func (m *mockAmazonKendraClient) Retrieve(ctx context.Context, params *kendra.RetrieveInput, optFns ...func(*kendra.Options)) (*kendra.RetrieveOutput, error) {
	return m.RetrieveOutput, m.RetrieveError
}

func (m *mockAmazonKendraClient) Query(ctx context.Context, params *kendra.QueryInput, optFns ...func(*kendra.Options)) (*kendra.QueryOutput, error) {
	return m.QueryOutput, m.QueryError
}
