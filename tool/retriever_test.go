package tool

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestRetriever(t *testing.T) {
	t.Parallel()

	t.Run("Run", func(t *testing.T) {
		t.Parallel()

		t.Run("Success", func(t *testing.T) {
			mockRetriever := &mockRetriever{
				docsResp: []schema.Document{
					{PageContent: "Document 1"},
					{PageContent: "Document 2"},
				},
			}
			retrieverTool := NewRetriever(mockRetriever, "Retriever", "A tool to retrieve documents")

			output, err := retrieverTool.Run(context.Background(), "query")
			assert.NoError(t, err)
			assert.Equal(t, "Document 1\n\nDocument 2", output)
		})

		t.Run("Error", func(t *testing.T) {
			mockRetriever := &mockRetriever{
				errorResp: errors.New("retriever error"),
			}
			retrieverTool := NewRetriever(mockRetriever, "Retriever", "A tool to retrieve documents")

			_, err := retrieverTool.Run(context.Background(), "query")
			assert.ErrorContains(t, err, "retriever error")
		})
	})

	t.Run("Getter", func(t *testing.T) {
		t.Parallel()

		mockRetriever := &mockRetriever{}
		retrieverTool := NewRetriever(mockRetriever, "Retriever", "A tool to retrieve documents")

		t.Run("Name", func(t *testing.T) {
			assert.Equal(t, "Retriever", retrieverTool.Name())
		})

		t.Run("Description", func(t *testing.T) {
			assert.Equal(t, "A tool to retrieve documents", retrieverTool.Description())
		})

		t.Run("ArgsType", func(t *testing.T) {
			assert.Equal(t, reflect.TypeOf(""), retrieverTool.ArgsType())
		})

		t.Run("Verbose", func(t *testing.T) {
			assert.False(t, retrieverTool.Verbose())
		})

		t.Run("Callbacks", func(t *testing.T) {
			assert.Equal(t, []schema.Callback(nil), retrieverTool.Callbacks())
		})
	})
}

// mockRetriever is a mock implementation of the schema.Retriever interface.
type mockRetriever struct {
	docsResp  []schema.Document
	errorResp error
}

// GetRelevantDocuments is the mock implementation of the GetRelevantDocuments method for mockRetriever.
func (m *mockRetriever) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	if m.errorResp != nil {
		return nil, m.errorResp
	}

	return m.docsResp, nil
}

// Callbacks is the mock implementation of the Callbacks method for mockRetriever.
func (m *mockRetriever) Callbacks() []schema.Callback {
	return nil
}

// Verbose is the mock implementation of the Verbose method for mockRetriever.
func (m *mockRetriever) Verbose() bool {
	return false
}
