package ai21

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCompletion(t *testing.T) {
	testCases := []struct {
		name          string
		model         string
		expectedError error
		expectedText  string
		mockClient    *mockHTTPClient
	}{
		{
			name:  "Success",
			model: "j2-light",
			mockClient: &mockHTTPClient{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader([]byte(`{
						"id": "test-id",
						"prompt": {
							"text": "Once upon a time"
						},
						"completions": [
							{
								"data": {
									"text": "in a land far, far away"
								},
								"finishReason": {
									"reason": "stop",
									"length": 27
								}
							}
						]
					}`))),
				},
				err:           nil,
				expectedURL:   "https://api.ai21.com/studio/v1/j2-light/complete",
				expectedToken: "Bearer test-api-key",
			},
			expectedError: nil,
			expectedText:  "in a land far, far away",
		},
		{
			name:          "UnknownModel",
			model:         "unknown-model",
			mockClient:    nil, // No mock client set, it should return an error.
			expectedError: fmt.Errorf("unknown model: unknown-model"),
			expectedText:  "",
		},
		{
			name:  "HTTPError",
			model: "j2-light",
			mockClient: &mockHTTPClient{
				err:           errors.New("HTTP error"),
				expectedURL:   "https://api.ai21.com/studio/v1/j2-light/complete",
				expectedToken: "Bearer test-api-key",
			},
			expectedError: errors.New("HTTP error"),
			expectedText:  "",
		},
	}

	for _, test := range testCases {
		client := New("test-api-key", func(o *Options) {
			o.HTTPClient = test.mockClient
		})

		res, err := client.CreateCompletion(context.TODO(), test.model, nil)
		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		}

		if test.expectedText != "" {
			assert.Equal(t, test.expectedText, res.Completions[0].Data.Text)
		}
	}
}

type mockHTTPClient struct {
	response      *http.Response
	err           error
	expectedURL   string
	expectedToken string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if req.URL.String() != m.expectedURL {
		return nil, fmt.Errorf("unexpected URL: %s", req.URL.String())
	}

	if req.Header.Get("Authorization") != m.expectedToken {
		return nil, errors.New("unauthorized")
	}

	if m.err != nil {
		return nil, m.err
	}

	return m.response, nil
}
