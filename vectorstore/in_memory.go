package vectorstore

import (
	"context"
	"sort"

	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/metric"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure InMemory satisfies the VectorStore interface.
var _ schema.VectorStore = (*InMemory)(nil)

// InMemoryItem represents an item stored in memory with its content, vector, and metadata.
type InMemoryItem struct {
	Content  string         `json:"content"`
	Vector   []float32      `json:"vector"`
	Metadata map[string]any `json:"metadata"`
}

// InMemoryOptions represents options for the in-memory vector store.
type InMemoryOptions struct {
	TopK int
}

// InMemory represents an in-memory vector store.
// Note: This implementation is intended for testing and demonstration purposes, not for production use.
type InMemory struct {
	embedder schema.Embedder
	data     []InMemoryItem
	opts     InMemoryOptions
}

// NewInMemory creates a new instance of the in-memory vector store.
func NewInMemory(embedder schema.Embedder, optFns ...func(*InMemoryOptions)) *InMemory {
	opts := InMemoryOptions{
		TopK: 3,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &InMemory{
		data:     make([]InMemoryItem, 0),
		embedder: embedder,
		opts:     opts,
	}
}

// AddDocuments adds a batch of documents to the InMemory vector store.
func (vs *InMemory) AddDocuments(ctx context.Context, docs []schema.Document) error {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.PageContent
	}

	vectors, err := vs.embedder.BatchEmbedText(ctx, texts)
	if err != nil {
		return err
	}

	for i, doc := range docs {
		vs.data = append(vs.data, InMemoryItem{
			Content:  doc.PageContent,
			Vector:   vectors[i],
			Metadata: doc.Metadata,
		})
	}

	return nil
}

// AddItem adds a single item to the InMemory vector store.
func (vs *InMemory) AddItem(item InMemoryItem) {
	vs.data = append(vs.data, item)
}

// Data returns the underlying data stored in the InMemory vector store.
func (vs *InMemory) Data() []InMemoryItem {
	return vs.data
}

// SimilaritySearch performs a similarity search with the given query in the InMemory vector store.
func (vs *InMemory) SimilaritySearch(ctx context.Context, query string) ([]schema.Document, error) {
	queryVector, err := vs.embedder.EmbedText(ctx, query)
	if err != nil {
		return nil, err
	}

	type searchResult struct {
		Item       InMemoryItem
		Similarity float32
	}

	results := make([]searchResult, len(vs.data))

	for i, item := range vs.data {
		similarity, err := metric.CosineSimilarity(queryVector, item.Vector)
		if err != nil {
			return nil, err
		}

		results[i] = searchResult{Item: item, Similarity: similarity}
	}

	// Sort results by similarity in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	docLen := util.Min(len(results), vs.opts.TopK)

	// Extract documents from sorted results
	documents := make([]schema.Document, docLen)
	for i := 0; i < docLen; i++ {
		documents[i] = schema.Document{
			PageContent: results[i].Item.Content,
			Metadata:    results[i].Item.Metadata,
		}
	}

	return documents, nil
}
