package retriever

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestMergeDocuments(t *testing.T) {
	retriever1 := &RetrieverMock{
		GetRelevantDocumentsFunc: func(ctx context.Context, query string) ([]schema.Document, error) {
			return []schema.Document{{PageContent: "Document 1"}}, nil
		},
	}

	retriever2 := &RetrieverMock{
		GetRelevantDocumentsFunc: func(ctx context.Context, query string) ([]schema.Document, error) {
			return []schema.Document{{PageContent: "Document 2"}, {PageContent: "Document 3"}}, nil
		},
	}

	retriever3 := &RetrieverMock{
		GetRelevantDocumentsFunc: func(ctx context.Context, query string) ([]schema.Document, error) {
			return []schema.Document{{PageContent: "Document 4"}, {PageContent: "Document 5"}}, nil
		},
	}

	// Mock the second retriever to return an empty result
	retriever4 := &RetrieverMock{
		GetRelevantDocumentsFunc: func(ctx context.Context, query string) ([]schema.Document, error) {
			return []schema.Document{}, nil
		},
	}

	t.Run("MergeDocuments returns merged documents from 2 retrievers", func(t *testing.T) {
		merger := NewMerger(retriever1, retriever2)

		query := "test query"
		expectedDocuments := []schema.Document{
			{PageContent: "Document 1"},
			{PageContent: "Document 2"},
			{PageContent: "Document 3"},
		}

		mergedDocuments, err := merger.GetRelevantDocuments(context.TODO(), query)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedDocuments, mergedDocuments)
	})

	t.Run("MergeDocuments returns merged documents from 3 retrievers", func(t *testing.T) {
		merger := NewMerger(retriever1, retriever2, retriever3)

		query := "test query"
		expectedDocuments := []schema.Document{
			{PageContent: "Document 1"},
			{PageContent: "Document 2"},
			{PageContent: "Document 3"},
			{PageContent: "Document 4"},
			{PageContent: "Document 5"},
		}

		mergedDocuments, err := merger.GetRelevantDocuments(context.TODO(), query)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedDocuments, mergedDocuments)
	})

	t.Run("MergeDocuments handles empty retriever results", func(t *testing.T) {
		merger := NewMerger(retriever1, retriever4)

		query := "test query"
		expectedDocuments := []schema.Document{{PageContent: "Document 1"}}

		// Mock the second retriever to return an empty result
		retriever2.GetRelevantDocumentsFunc = func(ctx context.Context, query string) ([]schema.Document, error) {
			return []schema.Document{}, nil
		}

		mergedDocuments, err := merger.GetRelevantDocuments(context.TODO(), query)
		assert.NoError(t, err)
		assert.Equal(t, expectedDocuments, mergedDocuments)
	})
}
