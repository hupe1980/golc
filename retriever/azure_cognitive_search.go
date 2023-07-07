package retriever

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure AzureCognitiveSearch satisfies the Retriever interface.
var _ schema.Retriever = (*AzureCognitiveSearch)(nil)

type AzureCognitiveSearchRequest struct {
	Search string `json:"search"`
	Top    uint   `json:"top"`
}

type AzureCognitiveSearchOptions struct {
	// Number of documents to query for
	TopK uint

	APIVersion string
	ContentKey string
	HTTPClient HTTPClient
}

type AzureCognitiveSearch struct {
	httpClient  HTTPClient
	apiKey      string
	serviceName string
	indexName   string
	opts        AzureCognitiveSearchOptions
}

func NewAzureCognitiveSearch(apiKey, serviceName, indexName string, optFns ...func(o *AzureCognitiveSearchOptions)) *AzureCognitiveSearch {
	opts := AzureCognitiveSearchOptions{
		TopK:       3,
		APIVersion: "2020-06-30",
		ContentKey: "content",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &AzureCognitiveSearch{
		httpClient: http.DefaultClient,
		opts:       opts,
	}
}

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

	items, ok := jsonMap["value"].([]map[string]any)
	if !ok {
		return nil, errors.New("bad response: value is missing")
	}

	docs := []schema.Document{}

	for _, item := range items {
		if content, ok := item[r.opts.ContentKey].(string); ok {
			docs = append(docs, schema.Document{
				PageContent: content,
				Metadata:    item,
			})
		}
	}

	return docs, nil
}

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

	res, err := r.httpClient.Do(httpReq)
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
