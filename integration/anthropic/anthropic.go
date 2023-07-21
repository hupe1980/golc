package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	HumanPrompt = "\n\nHuman:"
	AIPrompt    = "\n\nAssistant:"
)

type Options struct {
	APIUrl  string
	Version string
	SDK     string
}

type Client struct {
	httpClient *http.Client
	apiKey     string
	opts       Options
}

func New(apiKey string) *Client {
	opts := Options{
		APIUrl:  "https://api.anthropic.com",
		Version: "2023-01-01",
		SDK:     "golc-anthrophic-sdk",
	}

	return &Client{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
		opts:       opts,
	}
}

type CompletionRequest struct {
	Prompt      string            `json:"prompt"`
	Temperature float32           `json:"temperature,omitempty"`
	MaxTokens   int               `json:"max_tokens_to_sample"`
	Stop        []string          `json:"stop_sequences"`
	TopK        int               `json:"top_k,omitempty"`
	TopP        float32           `json:"top_p,omitempty"`
	Model       string            `json:"model"`
	Tags        map[string]string `json:"tags,omitempty"`
	Stream      bool              `json:"stream"`
}

type CompletionResponse struct {
	Completion string `json:"completion"`
	Stop       string `json:"stop"`
	StopReason string `json:"stop_reason"`
	Truncated  bool   `json:"truncated"`
	Exception  string `json:"exception"`
	LogID      string `json:"log_id"`
}

func (c *Client) Complete(ctx context.Context, request *CompletionRequest) (*CompletionResponse, error) {
	request.Stream = false

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.opts.APIUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-SDK", c.opts.SDK)
	req.Header.Set("Anthropic-Version", c.opts.Version)
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
