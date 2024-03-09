package ernie

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateChatCompletion(t *testing.T) {
	// Initialize the Ernie client with a mock HTTP client
	mockClient := &mockHTTPClient{}

	client := New("your-client-id", "your-client-secret", func(o *Options) {
		o.APIUrl = "https://example.com"
		o.HTTPClient = mockClient
	})

	t.Run("Successful Chat Completion", func(t *testing.T) {
		// Set up expected values
		expectedAccessToken := "your-access-token"
		expectedModel := "ernie-bot-3.5"
		expectedRequest := &ChatCompletionRequest{
			Messages: []Message{
				{
					Role:    "system",
					Content: "You are a helpful assistant.",
				},
			},
		}

		// Implement the DoFunc for the mock HTTP client
		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			// Assert the request
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "https://example.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions?access_token=your-access-token", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

			// Simulate a successful response
			response := &ChatCompletionResponse{
				Result: "A completed chat message.",
			}
			respBody, _ := json.Marshal(response)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(respBody)),
			}, nil
		}

		// Call CreateChatCompletion
		ctx := context.Background()
		client.accessToken = expectedAccessToken
		response, err := client.CreateChatCompletion(ctx, expectedModel, expectedRequest)

		// Assert the response and error
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "A completed chat message.", response.Result)
	})
}

func TestClient_CreateEmbedding(t *testing.T) {
	// Initialize the Ernie client with a mock HTTP client
	mockClient := &mockHTTPClient{}

	client := New("your-client-id", "your-client-secret", func(o *Options) {
		o.APIUrl = "https://example.com"
		o.HTTPClient = mockClient
	})

	t.Run("Successful Text Embedding", func(t *testing.T) {
		// Set up expected values
		expectedAccessToken := "your-access-token"
		expectedModel := "ernie-text-embedding"
		expectedRequest := EmbeddingRequest{
			Input: []string{"Text to embed."},
		}

		// Implement the DoFunc for the mock HTTP client
		mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			// Assert the request
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "https://example.com/rpc/2.0/ai_custom/v1/wenxinworkshop/embeddings/embedding-v1?access_token=your-access-token", req.URL.String())
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

			// Simulate a successful response
			response := &EmbeddingResponse{
				Data: []struct {
					Object    string    `json:"object"`
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Object:    "embedding",
						Embedding: []float32{0.1, 0.2, 0.3},
						Index:     0,
					},
				},
			}
			respBody, _ := json.Marshal(response)

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(respBody)),
			}, nil
		}

		// Call CreateEmbedding
		ctx := context.Background()
		client.accessToken = expectedAccessToken
		response, err := client.CreateEmbedding(ctx, expectedModel, expectedRequest)

		// Assert the response and error
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Data))
		assert.Equal(t, float32(0.1), response.Data[0].Embedding[0])
	})
}

// mockHTTPClient is a custom mock implementation of the HTTPClient interface.
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.Path, "oauth/2.0/token") {
		return nil, nil
	}

	if m.DoFunc != nil {
		return m.DoFunc(req)
	}

	return nil, nil
}
