package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	// HumanPrompt is the string used to indicate a human message in the conversation.
	HumanPrompt = "\n\nHuman:"

	// AIPrompt is the string used to indicate an assistant message in the conversation.
	AIPrompt = "\n\nAssistant:"
)

// Options represents the configuration options for the Anthropic client.
type Options struct {
	// The base URL of the Anthropic API.
	APIUrl string

	// The version of the Anthropic API to use.
	Version string

	// The SDK identifier used in the API requests.
	SDK string

	// The HTTP client to use for making API requests.
	HTTPClient HTTPClient
}

// Client represents the Anthropic API client.
type Client struct {
	apiKey string
	opts   Options
}

// New creates a new instance of the Anthropic API client with the given API key and optional configuration options.
func New(apiKey string, optFns ...func(o *Options)) *Client {
	opts := Options{
		APIUrl:     "https://api.anthropic.com",
		Version:    "2023-01-01",
		SDK:        "golc-anthrophic-sdk",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		apiKey: apiKey,
		opts:   opts,
	}
}

// CompletionRequest represents a request to the Anthropic API for text completion.
type CompletionRequest struct {
	// The input prompt for the completion.
	Prompt string `json:"prompt"`
	// The temperature for randomness in sampling.
	Temperature float32 `json:"temperature,omitempty"`
	// The maximum number of tokens to sample.
	MaxTokens int `json:"max_tokens_to_sample"`
	// List of strings to stop generation at.
	Stop []string `json:"stop_sequences"`
	// The number of highest probability tokens to use in sampling.
	TopK int `json:"top_k,omitempty"`
	// The cumulative probability for nucleus sampling.
	TopP float32 `json:"top_p,omitempty"`
	// The model to use for completion.
	Model string `json:"model"`
	// Additional tags for the completion.
	Tags map[string]string `json:"tags,omitempty"`
	// Flag to enable streaming response.
	Stream bool `json:"stream"`
}

// CompletionResponse represents the response from the Anthropic API for text completion.
type CompletionResponse struct {
	// The generated completion text.
	Completion string `json:"completion"`
	// The stop sequence that caused generation to stop.
	Stop string `json:"stop"`
	// The reason for stopping generation.
	StopReason string `json:"stop_reason"`
	// Flag indicating if the generated completion was truncated.
	Truncated bool `json:"truncated"`
	// The exception message if an error occurred during generation.
	Exception string `json:"exception"`
	// The log ID for the API request.
	LogID string `json:"log_id"`
}

// CreateCompletion sends a text completion request to the Anthropic API and returns the response.
func (c *Client) CreateCompletion(ctx context.Context, request *CompletionRequest) (*CompletionResponse, error) {
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

	resp, err := c.opts.HTTPClient.Do(req)
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
