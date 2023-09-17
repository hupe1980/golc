// Package ai21 provides a client for interacting with the ai21 text completion API.
package ai21

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hupe1980/golc/util"
)

// SupportedModels is a list of supported AI21 models.
var SupportedModels = []string{"j2-light", "j2-mid", "j2-ultra"}

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Options represents configuration options for the ai21 client.
type Options struct {
	// The base URL of the ai21 API.
	APIUrl string

	// The HTTP client to use for making API requests.
	HTTPClient HTTPClient
}

// Client is a client for interacting with the ai21 API.
type Client struct {
	apiKey string
	opts   Options
}

// New creates a new ai21 Client instance with the provided API key.
func New(apiKey string, optFns ...func(o *Options)) *Client {
	opts := Options{
		APIUrl:     "https://api.ai21.com/studio/v1",
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

// CompleteRequest represents a request for text completion.
type CompleteRequest struct {
	Prompt           string   `json:"prompt"`
	NumResults       int      `json:"numResults,omitempty"`
	MaxTokens        int      `json:"maxTokens,omitempty"`
	MinTokens        int      `json:"minTokens,omitempty"`
	Temperature      float64  `json:"temperature,omitempty"`
	TopP             float64  `json:"topP,omitempty"`
	StopSequences    []string `json:"stopSequences,omitempty"`
	TopKReturn       int      `json:"topKReturn,omitempty"`
	FrequencyPenalty Penalty  `json:"frequencyPenalty,omitempty"`
	PresencePenalty  Penalty  `json:"presencePenalty,omitempty"`
	CountPenalty     Penalty  `json:"countPenalty,omitempty"`
}

// Penalty represents penalty options for text completion.
type Penalty struct {
	Scale               int  `json:"scale"`
	ApplyToNumbers      bool `json:"applyToNumbers"`
	ApplyToPunctuations bool `json:"applyToPunctuations"`
	ApplyToStopwords    bool `json:"applyToStopwords"`
	ApplyToWhitespaces  bool `json:"applyToWhitespaces"`
	ApplyToEmojis       bool `json:"applyToEmojis"`
}

// CompleteResponse represents the response from a text completion request.
type CompleteResponse struct {
	ID          string       `json:"id"`
	Prompt      Prompt       `json:"prompt"`
	Completions []Completion `json:"completions"`
}

// Prompt represents the prompt text and tokens.
type Prompt struct {
	Text   string   `json:"text"`
	Tokens []Tokens `json:"tokens"`
}

// Tokens represents generated tokens and top tokens.
type Tokens struct {
	GeneratedToken GeneratedToken `json:"generatedToken"`
	TopTokens      interface{}    `json:"topTokens"`
	TextRange      TextRange      `json:"textRange"`
}

// GeneratedToken represents a generated token with log probabilities.
type GeneratedToken struct {
	Token      string  `json:"token"`
	Logprob    float64 `json:"logprob"`
	RawLogprob float64 `json:"raw_logprob"`
}

// TextRange represents a range of text in the prompt.
type TextRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Completion represents a completion result.
type Completion struct {
	Data         Data         `json:"data"`
	FinishReason FinishReason `json:"finishReason"`
}

// Data represents the completion text and tokens.
type Data struct {
	Text   string   `json:"text"`
	Tokens []Tokens `json:"tokens"`
}

// FinishReason represents the reason and length for completion finishing.
type FinishReason struct {
	Reason string `json:"reason"`
	Length int    `json:"length"`
}

// CreateCompletion sends a text completion request to the ai21 API and returns the response.
func (c *Client) CreateCompletion(ctx context.Context, model string, req *CompleteRequest) (*CompleteResponse, error) {
	if !util.Contains(SupportedModels, model) {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	url := fmt.Sprintf("%s/%s/complete", c.opts.APIUrl, model)

	body, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}

	completion := CompleteResponse{}
	if err := json.Unmarshal(body, &completion); err != nil {
		return nil, err
	}

	return &completion, nil
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
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

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
		return nil, fmt.Errorf("ai21 API returned unexpected status code: %d", res.StatusCode)
	}

	return resBody, nil
}
