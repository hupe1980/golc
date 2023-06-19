package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RestClient struct {
	apiKey string
	target string
}

func NewRestClient(apiKey string, endpoint Endpoint) (*RestClient, error) {
	target := endpoint.String()

	return &RestClient{
		apiKey: apiKey,
		target: target,
	}, nil
}

func (p *RestClient) Upsert(ctx context.Context, req *UpsertRequest) (*UpsertResponse, error) {
	reqURL := fmt.Sprintf("https://%s/vectors/upsert", p.target)

	res, err := p.doRequest(ctx, http.MethodPost, reqURL, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	upsertResponse := UpsertResponse{}
	if err := json.Unmarshal(body, &upsertResponse); err != nil {
		return nil, err
	}

	return &upsertResponse, nil
}

func (p *RestClient) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	params := make(url.Values)

	for _, id := range req.IDs {
		params.Add("ids", id)
	}

	if req.Namespace != "" {
		params.Add("format", "json")
	}

	reqURL := fmt.Sprintf("https://%s/vectors/fetch?%s", p.target, params.Encode())

	res, err := p.doRequest(ctx, http.MethodPost, reqURL, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fetchResponse := FetchResponse{}
	if err := json.Unmarshal(body, &fetchResponse); err != nil {
		return nil, err
	}

	return &fetchResponse, nil
}

func (p *RestClient) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	reqURL := fmt.Sprintf("https://%s/query", p.target)

	res, err := p.doRequest(ctx, http.MethodPost, reqURL, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	queryResponse := QueryResponse{}
	if err := json.Unmarshal(body, &queryResponse); err != nil {
		return nil, err
	}

	return &queryResponse, nil
}

func (p *RestClient) Close() error {
	return nil
}

func (p *RestClient) doRequest(ctx context.Context, method string, url string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("Api-Key", p.apiKey)

	return http.DefaultClient.Do(httpReq)
}
