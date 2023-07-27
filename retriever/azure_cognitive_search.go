package retriever

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure AzureCognitiveSearch satisfies the Retriever interface.
var _ schema.Retriever = (*AzureCognitiveSearch)(nil)

// AzureCognitiveSearchRequest represents the request payload for Azure Cognitive Search.
type AzureCognitiveSearchRequest struct {
	Search string `json:"search"`
	Top    uint   `json:"top"`
}

// AzureCognitiveSearchOptions contains options for configuring the AzureCognitiveSearch retriever.
type AzureCognitiveSearchOptions struct {
	*schema.CallbackOptions
	// Number of documents to query for
	TopK uint

	// Azure Cognitive Search API version.
	APIVersion string

	// Key to extract content from response.
	ContentKey string

	// HTTP client to use for making requests.
	HTTPClient HTTPClient
}

// AzureCognitiveSearch is a retriever implementation for Azure Cognitive Search service.
type AzureCognitiveSearch struct {
	apiKey      string
	serviceName string
	indexName   string
	opts        AzureCognitiveSearchOptions
}

// NewAzureCognitiveSearch creates a new instance of AzureCognitiveSearch retriever with the provided options.
func NewAzureCognitiveSearch(apiKey, serviceName, indexName string, optFns ...func(o *AzureCognitiveSearchOptions)) *AzureCognitiveSearch {
	opts := AzureCognitiveSearchOptions{
		TopK:       3,
		APIVersion: "2020-06-30",
		ContentKey: "content",
		HTTPClient: http.DefaultClient,
		CallbackOptions: &schema.CallbackOptions{
			Verbose: golc.Verbose,
		},
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AzureCognitiveSearch{
		apiKey:      apiKey,
		serviceName: serviceName,
		indexName:   indexName,
		opts:        opts,
	}
}

// GetRelevantDocuments retrieves relevant documents for the given query using Azure Cognitive Search.
func (r *AzureCognitiveSearch) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	url := fmt.Sprintf("https://%s.search.windows.net/indexes/%s/docs/search?api-version=%s", r.serviceName, r.indexName, r.opts.APIVersion)

	body, err := r.doRequest(ctx, http.MethodPost, url, &AzureCognitiveSearchRequest{
		Search: query,
		Top:    r.opts.TopK,
	})
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]any)
	if err := json.Unmarshal(body, &jsonMap); err != nil {
		return nil, err
	}

	items, ok := jsonMap["value"].([]any)
	if !ok {
		return nil, errors.New("bad response: value is missing")
	}

	docs := []schema.Document{}

	for _, item := range items {
		itemMap, _ := item.(map[string]any)

		if content, ok := itemMap[r.opts.ContentKey]; ok {
			docs = append(docs, schema.Document{
				PageContent: content.(string),
				Metadata:    itemMap,
			})
		}
	}

	return docs, nil
}

// Verbose returns the verbosity setting of the retriever.
func (r *AzureCognitiveSearch) Verbose() bool {
	return r.opts.CallbackOptions.Verbose
}

// Callbacks returns the registered callbacks of the retriever.
func (r *AzureCognitiveSearch) Callbacks() []schema.Callback {
	return r.opts.CallbackOptions.Callbacks
}

// doRequest sends an HTTP request to the Azure Cognitive Search service.
func (r *AzureCognitiveSearch) doRequest(ctx context.Context, method string, url string, payload any) ([]byte, error) {
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
	httpReq.Header.Set("Api-Key", r.apiKey)

	res, err := r.opts.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure cognitive search error: %s", string(resBody))
	}

	return resBody, nil
}
