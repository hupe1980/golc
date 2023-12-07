package retriever

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestBedrockKnowledgeBases_GetRelevantDocuments(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		retrieveOutput *bedrockagentruntime.RetrieveOutput
		retrieveError  error
		expectedDocs   []schema.Document
		expectedError  error
	}{
		{
			name: "Retrieve success",
			retrieveOutput: &bedrockagentruntime.RetrieveOutput{
				RetrievalResults: []types.KnowledgeBaseRetrievalResult{
					{
						Content: &types.RetrievalResultContent{
							Text: aws.String("Content 1"),
						},
						Location: &types.RetrievalResultLocation{
							Type: types.RetrievalResultLocationTypeS3,
							S3Location: &types.RetrievalResultS3Location{
								Uri: aws.String("URI 1"),
							},
						},
						Score: aws.Float64(0.9),
					},
					{
						Content: &types.RetrievalResultContent{
							Text: aws.String("Content 2"),
						},
						Location: &types.RetrievalResultLocation{
							Type: types.RetrievalResultLocationTypeS3,
							S3Location: &types.RetrievalResultS3Location{
								Uri: aws.String("URI 2"),
							},
						},
						Score: aws.Float64(0.8),
					},
				},
			},
			retrieveError: nil,
			expectedDocs: []schema.Document{
				{
					PageContent: "Content 1",
					Metadata: map[string]any{
						"location": "URI 1",
						"score":    float64(0.9),
					},
				},
				{
					PageContent: "Content 2",
					Metadata: map[string]any{
						"location": "URI 2",
						"score":    float64(0.8),
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Retrieve error",
			retrieveOutput: nil,
			retrieveError:  errors.New("retrieve error"),
			expectedDocs:   nil,
			expectedError:  errors.New("retrieve error"),
		},
		{
			name:           "No retrieval results",
			retrieveOutput: &bedrockagentruntime.RetrieveOutput{RetrievalResults: []types.KnowledgeBaseRetrievalResult{}},
			retrieveError:  nil,
			expectedDocs:   []schema.Document{},
			expectedError:  nil,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the retriever with the mock client
			r := NewBedrockKnowledgeBases(&mockBedrockAgentRuntimeClient{
				RetrieveOutput: tt.retrieveOutput,
				RetrieveError:  tt.retrieveError,
			}, "knowledge-base-id")

			// Call the method under test
			docs, err := r.GetRelevantDocuments(context.Background(), "query")

			// Check the results
			assert.Equal(t, tt.expectedDocs, docs)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// mockBedrockAgentRuntimeClient is a mock implementation of the BedrockAgentRuntimeClient interface.
type mockBedrockAgentRuntimeClient struct {
	RetrieveOutput *bedrockagentruntime.RetrieveOutput
	RetrieveError  error
}

func (m *mockBedrockAgentRuntimeClient) Retrieve(ctx context.Context, params *bedrockagentruntime.RetrieveInput, optFns ...func(*bedrockagentruntime.Options)) (*bedrockagentruntime.RetrieveOutput, error) {
	return m.RetrieveOutput, m.RetrieveError
}
