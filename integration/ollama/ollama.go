package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hupe1980/golc/integration/stream"
)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOptions struct {
	// The HTTP client to use for making API requests.
	HTTPClient HTTPClient
}

type Client struct {
	apiURL string
	opts   ClientOptions
}

func New(apiURL string, optFns ...func(o *ClientOptions)) *Client {
	opts := ClientOptions{
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		apiURL: apiURL,
		opts:   opts,
	}
}

func (c *Client) CreateGeneration(ctx context.Context, req *GenerationRequest) (*GenerationResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/generate", c.apiURL), req)
	if err != nil {
		return nil, err
	}

	completion := GenerationResponse{}
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, err
	}

	return &completion, nil
}

func (c *Client) CreateGenerationStream(ctx context.Context, req *GenerationRequest) (*GenerationStream, error) {
	res, err := c.doStreamRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/generate", c.apiURL), req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return nil, err
	}

	return &GenerationStream{
		Stream: stream.NewStream[GenerationResponse](res),
	}, nil
}

func (c *Client) CreateChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", c.apiURL), req)
	if err != nil {
		return nil, err
	}

	completion := ChatResponse{}
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, err
	}

	return &completion, nil
}

func (c *Client) CreateChatStream(ctx context.Context, req *ChatRequest) (*ChatStream, error) {
	res, err := c.doStreamRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/chat", c.apiURL), req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return nil, err
	}

	return &ChatStream{
		Stream: stream.NewStream[ChatResponse](res),
	}, nil
}

func (c *Client) CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/api/embeddings", c.apiURL), req)
	if err != nil {
		return nil, err
	}

	embedding := EmbeddingResponse{}
	if err := json.Unmarshal(body, &embedding); err != nil {
		return nil, err
	}

	return &embedding, nil
}

// doRequest sends an HTTP request to the specified URL with the given method and payload.
func (c *Client) doRequest(ctx context.Context, method string, url string, payload any) ([]byte, error) {
	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := c.opts.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		if err := json.Unmarshal(resBody, &errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("ollama API error: %s", errorResponse.Message)
	}

	return resBody, nil
}

func (c *Client) doStreamRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
	var buf *bytes.Buffer

	if data != nil {
		bts, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(bts)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")

	res, err := c.opts.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		errorResponse := ErrorResponse{}
		if err := json.Unmarshal(resBody, &errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("ollama API error: %s", errorResponse.Message)
	}

	return res, nil
}
