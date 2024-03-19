package vectorstore

import (
	"container/heap"
	"context"

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

// priorityQueueItem represents an item in the priority queue.
type priorityQueueItem struct {
	Data     InMemoryItem // Data associated with the item
	distance float32      // Distance from the query vector

	index int // Index of the item in the priority queue
}

// priorityQueue is a priority queue for InMemoryItem items.
type priorityQueue []*priorityQueueItem

// Len returns the length of the priority queue.
func (pq priorityQueue) Len() int { return len(pq) }

// Less reports whether the element with index i should sort before the element with index j.
func (pq priorityQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }

// Swap swaps the elements with indexes i and j.
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item to the priority queue.
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item, _ := x.(*priorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority (distance).
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]

	return item
}

// Top returns the top element of the priority queue.
func (pq *priorityQueue) Top() any {
	return (*pq)[0]
}

// DistanceFunc represents a function for calculating the distance between two vectors
type DistanceFunc func(v1, v2 []float32) (float32, error)

// InMemoryOptions represents options for the in-memory vector store.
type InMemoryOptions struct {
	TopK         int
	DistanceFunc DistanceFunc
}

// InMemory represents an in-memory vector store.
type InMemory struct {
	embedder schema.Embedder
	data     []InMemoryItem
	opts     InMemoryOptions
}

// NewInMemory creates a new instance of the in-memory vector store.
func NewInMemory(embedder schema.Embedder, optFns ...func(*InMemoryOptions)) *InMemory {
	opts := InMemoryOptions{
		TopK:         3,
		DistanceFunc: metric.SquaredL2,
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

	topCandidates := &priorityQueue{}

	heap.Init(topCandidates)

	type searchResult struct {
		Item       InMemoryItem
		Similarity float32
	}

	results := make([]searchResult, len(vs.data))

	for _, item := range vs.data {
		similarity, err := vs.opts.DistanceFunc(queryVector, item.Vector)
		if err != nil {
			return nil, err
		}

		if topCandidates.Len() < vs.opts.TopK {
			heap.Push(topCandidates, &priorityQueueItem{
				Data:     item,
				distance: similarity,
			})

			continue
		}

		largestDist, _ := topCandidates.Top().(*priorityQueueItem)

		if similarity < largestDist.distance {
			_ = heap.Pop(topCandidates)

			heap.Push(topCandidates, &priorityQueueItem{
				Data:     item,
				distance: similarity,
			})
		}
	}

	docLen := util.Min(len(results), vs.opts.TopK)

	// Extract documents from sorted results
	documents := make([]schema.Document, 0, docLen)

	for topCandidates.Len() > 0 {
		item, _ := heap.Pop(topCandidates).(*priorityQueueItem)

		documents = append(documents, schema.Document{
			PageContent: item.Data.Content,
			Metadata:    item.Data.Metadata,
		})
	}

	return documents, nil
}
