package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("CreateCompletion", func(t *testing.T) {
		mockResponse := CompletionResponse{
			Completion: "Hello, world!",
			Stop:       "",
			StopReason: "",
			Truncated:  false,
			Exception:  "",
			LogID:      "12345",
		}
		mockPayload, err := json.Marshal(mockResponse)
		assert.NoError(t, err)

		request := &CompletionRequest{
			Prompt:      "Hello, Anthropic!",
			Temperature: 0.8,
			MaxTokens:   50,
			Stop:        []string{"\n\nHuman:", "\n\nAssistant:"},
			TopK:        5,
			TopP:        0.7,
			Model:       "gpt-3.5-turbo",
		}

		// Create the mock HTTP client and set the Do function to return the mockResponse.
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", req.Header.Get("Accept"))
				assert.Equal(t, "golc-anthrophic-sdk", req.Header.Get("Anthropic-SDK"))
				assert.Equal(t, "2023-01-01", req.Header.Get("Anthropic-Version"))
				assert.Equal(t, "api-key", req.Header.Get("X-API-Key"))

				body, bErr := io.ReadAll(req.Body)
				assert.NoError(t, bErr)

				defer req.Body.Close()

				var request CompletionRequest
				err = json.Unmarshal(body, &request)
				assert.NoError(t, err)

				assert.Equal(t, "Hello, Anthropic!", request.Prompt)
				assert.Equal(t, float32(0.8), request.Temperature)
				assert.Equal(t, 50, request.MaxTokens)
				assert.Equal(t, []string{"\n\nHuman:", "\n\nAssistant:"}, request.Stop)
				assert.Equal(t, 5, request.TopK)
				assert.Equal(t, float32(0.7), request.TopP)
				assert.Equal(t, "gpt-3.5-turbo", request.Model)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(mockPayload)),
				}, nil
			},
		}

		client := New("api-key", func(o *Options) {
			o.HTTPClient = mockClient
		})

		response, err := client.CreateCompletion(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, &mockResponse, response)
	})

	t.Run("CreateCompletion_Error", func(t *testing.T) {
		request := &CompletionRequest{
			Prompt: "Hello, Anthropic!",
		}

		// Create the mock HTTP client and set the Do function to return an error.
		mockClient := &mockHTTPClient{
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, assert.AnError
			},
		}

		client := New("api-key", func(o *Options) {
			o.HTTPClient = mockClient
		})

		response, err := client.CreateCompletion(context.Background(), request)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

// mockHTTPClient is a mock implementation of the HTTPClient interface.
type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}
