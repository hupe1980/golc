package chain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	t.Run("Valid Input", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			if strings.HasSuffix(prompt, "API url:") {
				return "https://galaxy.org"
			}

			return "42"
		})

		api, err := NewAPI(fake, "doc", func(o *APIOptions) {
			o.HTTPClient = &mockHTTPClient{
				Response:   "number 42",
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			}
		})
		require.NoError(t, err)

		answer, err := golc.SimpleCall(context.Background(), api, "What is the Ultimate Answer to the question of Life, the Universe, and Everything?")
		require.NoError(t, err)
		require.Equal(t, "42", answer)
	})

	t.Run("Invalid Input Key", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			if strings.HasSuffix(prompt, "API url:") {
				return "https://galaxy.org"
			}
			return "42"
		})

		api, err := NewAPI(fake, "doc", func(o *APIOptions) {
			o.HTTPClient = &mockHTTPClient{
				Response:   "number 42",
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			}
		})
		require.NoError(t, err)

		_, err = golc.Call(context.Background(), api, schema.ChainValues{"invalid_key": "What is the Ultimate Answer?"})
		require.Error(t, err)
		require.EqualError(t, fmt.Errorf("invalid input values: no value for inputKey %s", api.InputKeys()[0]), err.Error())
	})

	t.Run("Invalid API URL", func(t *testing.T) {
		fake := llm.NewFake(func(prompt string) string {
			if strings.HasSuffix(prompt, "API url:") {
				return "https://galaxy.org"
			}
			return "42"
		})

		api, err := NewAPI(fake, "doc", func(o *APIOptions) {
			o.HTTPClient = &mockHTTPClient{
				Response:   "number 42",
				Status:     "200 OK",
				StatusCode: http.StatusOK,
			}
			o.VerifyURL = func(url string) bool {
				return false
			}
		})
		require.NoError(t, err)

		_, err = golc.SimpleCall(context.Background(), api, "What is the Ultimate Answer?")
		require.Error(t, err)
		require.EqualError(t, errors.New("invalid API URL: https://galaxy.org"), err.Error())
	})
}

type mockHTTPClient struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // http.StatusOK
	Response   string
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     c.Status,
		StatusCode: c.StatusCode,
		Body:       io.NopCloser(strings.NewReader(c.Response)),
	}, nil
}
