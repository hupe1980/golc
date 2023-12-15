// Package ernie provides a client library for interacting with the Ernie API, which offers natural language processing (NLP) capabilities, including chat completion and text embedding. Ernie is designed to assist with tasks such as generating human-like text responses or obtaining embeddings for text data.
package ernie

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// chatModelSuffixMap maps model names to their corresponding API endpoints for chat completion.
var chatModelSuffixMap = map[string]string{
	"ernie-bot-3.5":   "completions",
	"ernie-bot-turbo": "eb-instant",
}

// embeddingModelSuffixMap maps model names to their corresponding API endpoints for text embedding.
var embeddingModelSuffixMap = map[string]string{
	"ernie-text-embedding": "embedding-v1",
}

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Options represents configuration options for the Ernie client.
type Options struct {
	// The base URL of the Ernie API.
	APIUrl string

	// The HTTP client to use for making API requests.
	HTTPClient HTTPClient
}

// Client represents a client for interacting with the Ernie API.
type Client struct {
	clientID     string
	clientSecret string
	opts         Options
	accessToken  string
	mu           sync.RWMutex
}

// New creates a new instance of the Ernie client.
func New(clientID, clientSecret string, optFns ...func(o *Options)) *Client {
	opts := Options{
		APIUrl:     "https://aip.baidubce.com",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		opts:         opts,
	}
}

// Message represents a chat message with role and content.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request for chat completion.
type ChatCompletionRequest struct {
	Messages     []Message `json:"messages"`
	Temperature  float64   `json:"temperature,omitempty"`
	TopP         float64   `json:"top_p,omitempty"`
	PenaltyScore float64   `json:"penalty_score,omitempty"`
	Stream       bool      `json:"stream,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
}

// ChatCompletionResponse represents the response from chat completion API.
type ChatCompletionResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	SentenceID       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// CreateChatCompletion generates chat completion using the specified model and request.
func (c *Client) CreateChatCompletion(ctx context.Context, model string, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if c.accessToken == "" {
		err := c.requestAccessToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	chatCompletion, err := c.doChatCompletionRequest(ctx, model, request)
	if err != nil {
		return nil, err
	}

	if chatCompletion.ErrorCode == 111 { // access_token expired, refresh it
		err := c.requestAccessToken(ctx)
		if err != nil {
			return nil, err
		}

		return c.doChatCompletionRequest(ctx, model, request)
	}

	return chatCompletion, nil
}

// EmbeddingRequest represents a request for text embedding.
type EmbeddingRequest struct {
	Input []string `json:"input"`
}

// EmbeddingResponse represents the response from text embedding API.
type EmbeddingResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Data    []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// CreateEmbedding generates text embedding using the specified model and request.
func (c *Client) CreateEmbedding(ctx context.Context, model string, request EmbeddingRequest) (*EmbeddingResponse, error) {
	if c.accessToken == "" {
		err := c.requestAccessToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	embedding, err := c.doCreateEmbedding(ctx, model, request)
	if err != nil {
		return nil, err
	}

	if embedding.ErrorCode == 111 { // access_token expired, refresh it
		err := c.requestAccessToken(ctx)
		if err != nil {
			return nil, err
		}

		return c.doCreateEmbedding(ctx, model, request)
	}

	return embedding, nil
}

func (c *Client) doCreateEmbedding(ctx context.Context, model string, request EmbeddingRequest) (*EmbeddingResponse, error) {
	suffix, ok := embeddingModelSuffixMap[model]
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	params := make(url.Values)
	params.Add("access_token", c.accessToken)

	url := fmt.Sprintf("%s/rpc/2.0/ai_custom/v1/wenxinworkshop/embeddings/%s?%s", c.opts.APIUrl, suffix, params.Encode())

	res, err := c.doRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, err
	}

	embedding := EmbeddingResponse{}
	if err := json.Unmarshal(res, &embedding); err != nil {
		return nil, err
	}

	return &embedding, nil
}

func (c *Client) doChatCompletionRequest(ctx context.Context, model string, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	suffix, ok := chatModelSuffixMap[model]
	if !ok {
		return nil, fmt.Errorf("unknown model: %s", model)
	}

	params := make(url.Values)
	params.Add("access_token", c.accessToken)

	url := fmt.Sprintf("%s/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/%s?%s", c.opts.APIUrl, suffix, params.Encode())

	res, err := c.doRequest(ctx, http.MethodPost, url, request)
	if err != nil {
		return nil, err
	}

	chatCompletion := ChatCompletionResponse{}
	if err := json.Unmarshal(res, &chatCompletion); err != nil {
		return nil, err
	}

	return &chatCompletion, nil
}

// authResponse represents the response from the authentication API.
// see https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Ilkkrb0i5
type authResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// requestAccessToken requests a new access token for authentication.
func (c *Client) requestAccessToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	params := make(url.Values)
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", c.clientID)
	params.Add("client_secret", c.clientSecret)

	url := fmt.Sprintf("%s/oauth/2.0/token?%s", c.opts.APIUrl, params.Encode())

	res, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	auth := authResponse{}
	if err := json.Unmarshal(res, &auth); err != nil {
		return err
	}

	if auth.Error != "" {
		return fmt.Errorf("ernie bot error: %s", auth.Error)
	}

	c.accessToken = auth.AccessToken

	return nil
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
		return nil, fmt.Errorf("completion API returned unexpected status code: %d", res.StatusCode)
	}

	return resBody, nil
}
