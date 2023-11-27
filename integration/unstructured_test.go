package integration

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnstructuredPartition(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)

	defer os.Remove(file.Name())

	// Create a sample response JSON
	responseJSON := `[{"type":"NarrativeText","element_id":"mock_element_id","metadata":{"filetype":"application/pdf","languages":["eng"],"page_number":1,"filename":"testfile.pdf"},"text":"Mock text"}]`

	// Create a mock HTTP client with a predefined response
	mockUnstructuredHTTPClient := &mockUnstructuredHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Contains(t, req.Header.Get("Content-Type"), "multipart/form-data")
			assert.Equal(t, "mock_api_key", req.Header.Get("unstructured-api-key"))

			// Simulate a successful response
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(responseJSON)),
			}, nil
		},
	}

	// Create an instance of Unstructured with the mock HTTP client
	unstructuredClient := NewUnstructured("mock_api_key", func(o *UnstructuredOptions) {
		o.HTTPClient = mockUnstructuredHTTPClient
	})

	// Create a test case
	t.Run("Partition", func(t *testing.T) {
		// Call the Partition method with the mock file
		output, err := unstructuredClient.Partition(context.Background(), &PartitionInput{File: file})
		assert.NoError(t, err)

		// Assert the expected output
		expectedOutput := []PartitionOutput{
			{
				Type:      "NarrativeText",
				ElementID: "mock_element_id",
				Metadata: struct {
					Filetype   string   `json:"filetype"`
					Languages  []string `json:"languages"`
					PageNumber int      `json:"page_number"`
					Filename   string   `json:"filename"`
				}{
					Filetype:   "application/pdf",
					Languages:  []string{"eng"},
					PageNumber: 1,
					Filename:   "testfile.pdf",
				},
				Text: "Mock text",
			},
		}
		assert.Equal(t, expectedOutput, output)
	})
}

// mockUnstructuredHTTPClient is a custom mock for the HTTP client.
type mockUnstructuredHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do is the implementation of the Do method for the mock.
func (m *mockUnstructuredHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}

	return nil, errors.New("mockUnstructuredHTTPClient: DoFunc not set")
}
