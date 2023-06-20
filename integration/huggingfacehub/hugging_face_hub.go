package huggingfacehub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiEndpoint = "https://api-inference.huggingface.co"

type Client struct {
	apiToken string
	repoID   string
	task     string
}

func New(apiToken, repoID, task string) *Client {
	return &Client{
		apiToken: apiToken,
		repoID:   repoID,
		task:     task,
	}
}

func (hf *Client) Summarization(ctx context.Context, req *SummarizationRequest) (*SummarizationResponse, error) {
	reqURL := fmt.Sprintf("%s/pipeline/%s/%s", apiEndpoint, hf.task, hf.repoID)

	res, err := hf.doRequest(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("hugging faces error: %s", errResp.Error)
	}

	summarizationResponse := SummarizationResponse{}
	if err := json.Unmarshal(body, &summarizationResponse); err != nil {
		return nil, err
	}

	return &summarizationResponse, nil
}

func (hf *Client) TextGeneration(ctx context.Context, req *TextGenerationRequest) (TextGenerationResponse, error) {
	reqURL := fmt.Sprintf("%s/pipeline/%s/%s", apiEndpoint, hf.task, hf.repoID)

	res, err := hf.doRequest(ctx, http.MethodPost, reqURL, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("hugging faces error: %s", errResp.Error)
	}

	textGenerations := TextGenerationResponse{}
	if err := json.Unmarshal(body, &textGenerations); err != nil {
		return nil, err
	}

	return textGenerations, nil
}

func (hf *Client) Text2TextGeneration(ctx context.Context, req *Text2TextGenerationRequest) (Text2TextGenerationResponse, error) {
	reqURL := fmt.Sprintf("%s/pipeline/%s/%s", apiEndpoint, hf.task, hf.repoID)

	res, err := hf.doRequest(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("hugging faces error: %s", errResp.Error)
	}

	text2TextGenerationResponse := Text2TextGenerationResponse{}
	if err := json.Unmarshal(body, &text2TextGenerationResponse); err != nil {
		return nil, err
	}

	return text2TextGenerationResponse, nil
}

func (hf *Client) doRequest(ctx context.Context, method string, url string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hf.apiToken))

	return http.DefaultClient.Do(httpReq)
}
