package retriever

import (
	"context"
	"errors"
	"net/http"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure retrieverMock satisfies the Retriever interface.
var _ schema.Retriever = (*retrieverMock)(nil)

type retrieverMock struct {
	GetRelevantDocumentsFunc func(ctx context.Context, query string) ([]schema.Document, error)
}

func (m *retrieverMock) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	if m.GetRelevantDocumentsFunc != nil {
		return m.GetRelevantDocumentsFunc(ctx, query)
	}

	return nil, nil
}

func (m *retrieverMock) Verbose() bool {
	return false
}

func (m *retrieverMock) Callbacks() []schema.Callback {
	return nil
}

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.doFunc != nil {
		return c.doFunc(req)
	}

	return nil, errors.New("mock DoFunc is not set")
}
