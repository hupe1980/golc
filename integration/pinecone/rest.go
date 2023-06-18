package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(httpReq)
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
