package retriever

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAzureCognitiveSearch(t *testing.T) {
	t.Run("GetRelevantDocuments", func(t *testing.T) {
		// Create a mock HTTP client with expected response
		mockResp := `{"value": [{"content": "Document 1"}, {"content": "Document 2"}]}`
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
				}
				return resp, nil
			},
		}

		// Create an instance of AzureCognitiveSearch with the mock HTTP client
		retriever := NewAzureCognitiveSearch("apiKey", "serviceName", "indexName",
			func(o *AzureCognitiveSearchOptions) {
				o.HTTPClient = mockClient
			},
		)

		// Call the GetRelevantDocuments method
		docs, err := retriever.GetRelevantDocuments(context.Background(), "query")

		// Assert that the documents and error are as expected
		assert.NoError(t, err)
		assert.Len(t, docs, 2)
		assert.Equal(t, "Document 1", docs[0].PageContent)
		assert.Equal(t, "Document 2", docs[1].PageContent)
	})

	t.Run("GetRelevantDocuments_Error", func(t *testing.T) {
		// Create a mock HTTP client that returns an error
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("HTTP client error")
			},
		}

		// Create an instance of AzureCognitiveSearch with the mock HTTP client
		retriever := NewAzureCognitiveSearch("apiKey", "serviceName", "indexName",
			func(o *AzureCognitiveSearchOptions) {
				o.HTTPClient = mockClient
			},
		)

		// Call the GetRelevantDocuments method
		docs, err := retriever.GetRelevantDocuments(context.Background(), "query")

		// Assert that the error is as expected
		assert.Error(t, err)
		assert.Nil(t, docs)
	})

	t.Run("GetRelevantDocuments_InvalidResponse", func(t *testing.T) {
		// Create a mock HTTP client with invalid response JSON
		mockResp := `{"value": "invalid"}`
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
				}
				return resp, nil
			},
		}

		// Create an instance of AzureCognitiveSearch with the mock HTTP client
		retriever := NewAzureCognitiveSearch("apiKey", "serviceName", "indexName",
			func(o *AzureCognitiveSearchOptions) {
				o.HTTPClient = mockClient
			},
		)

		// Call the GetRelevantDocuments method
		docs, err := retriever.GetRelevantDocuments(context.Background(), "query")

		// Assert that the error is as expected
		assert.Error(t, err)
		assert.Nil(t, docs)
	})

	t.Run("GetRelevantDocuments_UnsuccessfulStatusCode", func(t *testing.T) {
		// Create a mock HTTP client with unsuccessful status code
		mockResp := `{"error": "Not found"}`
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
				}
				return resp, nil
			},
		}

		// Create an instance of AzureCognitiveSearch with the mock HTTP client
		retriever := NewAzureCognitiveSearch("apiKey", "serviceName", "indexName",
			func(o *AzureCognitiveSearchOptions) {
				o.HTTPClient = mockClient
			},
		)

		// Call the GetRelevantDocuments method
		docs, err := retriever.GetRelevantDocuments(context.Background(), "query")

		// Assert that the error is as expected
		assert.Error(t, err)
		assert.Nil(t, docs)
	})
}
