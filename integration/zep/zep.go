package zep

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	HTTPClient HTTPClient
	APIKey     string
	Version    string
}

type Client struct {
	baseURL string
	opts    Options
}

func New(baseURL string, optFns ...func(o *Options)) *Client {
	opts := Options{
		Version:    "v1",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		baseURL: baseURL,
		opts:    opts,
	}
}

// GetMemory retrieves memory for a specific session..
func (c *Client) GetMemory(ctx context.Context, sessionID string) (*Memory, error) {
	reqURL := fmt.Sprintf("%s/api/%s/sessions/%s/memory", c.baseURL, c.opts.Version, sessionID)

	body, err := c.doRequest(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	memory := Memory{}
	if err := json.Unmarshal(body, &memory); err != nil {
		return nil, err
	}

	return &memory, nil
}

// AddMemory adds a new memory to a specific session.
func (c *Client) AddMemory(ctx context.Context, sessionID string, memory *Memory) (string, error) {
	reqURL := fmt.Sprintf("%s/api/%s/sessions/%s/memory", c.baseURL, c.opts.Version, sessionID)

	body, err := c.doRequest(ctx, http.MethodPost, reqURL, memory)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// DeleteMemory deletes the memory of a specific session.
func (c *Client) DeleteMemory(ctx context.Context, sessionID string) (string, error) {
	reqURL := fmt.Sprintf("%s/api/%s/sessions/%s/memory", c.baseURL, c.opts.Version, sessionID)

	body, err := c.doRequest(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SearchMessages searches memory of a specific session based on search payload provided.
func (c *Client) SearchMessages(ctx context.Context, sessionID string, payload *SearchPayload) (*SearchResult, error) {
	reqURL := fmt.Sprintf("%s/api/%s/sessions/%s/search", c.baseURL, c.opts.Version, sessionID)

	body, err := c.doRequest(ctx, http.MethodPost, reqURL, payload)
	if err != nil {
		return nil, err
	}

	result := SearchResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

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

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	if c.opts.APIKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.opts.APIKey))
	}

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
		var apiErr APIError
		if err := json.Unmarshal(resBody, &apiErr); err != nil {
			return nil, fmt.Errorf("zep api error: %s", resBody)
		}

		return nil, fmt.Errorf("zep api error: %d - %s", apiErr.Code, apiErr.Message)
	}

	return resBody, nil
}
