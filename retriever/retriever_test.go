package retriever

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

type RetrieverMock struct {
	GetRelevantDocumentsFunc func(ctx context.Context, query string) ([]schema.Document, error)
}

func (m *RetrieverMock) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	if m.GetRelevantDocumentsFunc != nil {
		return m.GetRelevantDocumentsFunc(ctx, query)
	}

	return nil, nil
}
