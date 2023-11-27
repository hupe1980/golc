package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UnstructuredOptions represents options for configuring the Unstructured client.
type UnstructuredOptions struct {
	// BaseURL is the base URL of the Unstructured API.
	BaseURL string

	// HTTPClient is the HTTP client to use for making API requests.
	HTTPClient HTTPClient
}

// Unstructured is a client for interacting with the Unstructured API.
type Unstructured struct {
	apiKey string
	opts   UnstructuredOptions
}

// NewUnstructured creates a new Unstructured client with the provided API key.
func NewUnstructured(apiKey string, optFns ...func(o *UnstructuredOptions)) *Unstructured {
	opts := UnstructuredOptions{
		BaseURL:    "https://api.unstructured.io/general/v0/general",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Unstructured{
		apiKey: apiKey,
		opts:   opts,
	}
}

// PartitionInput represents the input for the Partition method.
type PartitionInput struct {
	File *os.File
}

// PartitionOutput represents the output of the Partition method.
type PartitionOutput struct {
	Type      string `json:"type"`
	ElementID string `json:"element_id"`
	Metadata  struct {
		Filetype   string   `json:"filetype"`
		Languages  []string `json:"languages"`
		PageNumber int      `json:"page_number"`
		Filename   string   `json:"filename"`
	} `json:"metadata"`
	Text string `json:"text"`
}

// Partition sends a file to the Unstructured API for partitioning and returns the partitioned content.
func (c *Unstructured) Partition(ctx context.Context, input *PartitionInput) ([]PartitionOutput, error) {
	fields := map[string]string{
		"strategy": "hi_res",
	}

	res, err := c.doMultipartRequest(ctx, c.opts.BaseURL, input.File, fields)
	if err != nil {
		return nil, err
	}

	output := []PartitionOutput{}
	if err := json.Unmarshal(res, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// doMultipartRequest performs a multipart request to the Unstructured API.
func (c *Unstructured) doMultipartRequest(ctx context.Context, url string, file *os.File, fields map[string]string) ([]byte, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("files", file.Name())
	if err != nil {
		return nil, err
	}

	if _, cErr := io.Copy(fw, file); cErr != nil {
		return nil, cErr
	}

	for k, v := range fields {
		if wErr := w.WriteField(k, v); wErr != nil {
			return nil, wErr
		}
	}

	// Close finishes the multipart message and writes the trailing boundary end line to the output.
	w.Close()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", w.FormDataContentType())
	httpReq.Header.Set("unstructured-api-key", c.apiKey)

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
		return nil, fmt.Errorf("unstructured API returned unexpected status code: %d", res.StatusCode)
	}

	return resBody, nil
}
