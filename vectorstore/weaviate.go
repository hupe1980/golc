package vectorstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

// Compile time check to ensure Weaviate satisfies the VectorStore interface.
var _ schema.VectorStore = (*Weaviate)(nil)

// WeaviateOptions contains options for configuring the Weaviate vector store.
type WeaviateOptions struct {
	// TextKey is the name of the property in the Weaviate objects where the text content is stored.
	TextKey string
	// TopK is the number of documents to retrieve in similarity search.
	TopK int
}

// Weaviate represents a Weaviate vector store.
type Weaviate struct {
	client    *weaviate.Client
	embedder  schema.Embedder
	indexName string
	opts      WeaviateOptions
}

// NewWeaviate creates a new Weaviate vector store with the given Weaviate client, embedder, index name, and optional configuration options.
func NewWeaviate(client *weaviate.Client, embedder schema.Embedder, indexName string, optFns ...func(*WeaviateOptions)) *Weaviate {
	opts := WeaviateOptions{
		TextKey: "text",
		TopK:    4,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Weaviate{
		client:    client,
		embedder:  embedder,
		indexName: indexName,
		opts:      opts,
	}
}

// AddDocuments adds a batch of documents to the Weaviate vector store.
func (vs *Weaviate) AddDocuments(ctx context.Context, docs []schema.Document) error {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.PageContent
	}

	vectors, err := vs.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return err
	}

	objects := make([]*models.Object, 0, len(docs))

	for i, doc := range docs {
		metadata := make(map[string]any, len(doc.Metadata))
		for key, value := range doc.Metadata {
			metadata[key] = value
		}

		metadata[vs.opts.TextKey] = doc.PageContent

		objects = append(objects, &models.Object{
			Class:      vs.indexName,
			ID:         strfmt.UUID(uuid.New().String()),
			Vector:     float64ToFloat32(vectors[i]),
			Properties: metadata,
		})
	}

	if _, err := vs.client.Batch().ObjectsBatcher().WithObjects(objects...).Do(ctx); err != nil {
		return err
	}

	return nil
}

// SimilaritySearch performs a similarity search with the given query in the Weaviate vector store.
func (vs *Weaviate) SimilaritySearch(ctx context.Context, query string) ([]schema.Document, error) {
	vector, err := vs.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err := vs.client.GraphQL().Get().
		WithNearVector(vs.client.GraphQL().NearVectorArgBuilder().WithVector(float64ToFloat32(vector))).
		WithClassName(vs.indexName).
		WithLimit(vs.opts.TopK).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if len(res.Errors) > 0 {
		messages := make([]string, len(res.Errors))
		for i, e := range res.Errors {
			messages[i] = e.Message
		}

		return nil, fmt.Errorf("weaviate errors: %s", strings.Join(messages, ", "))
	}

	data, ok := res.Data["Get"].(map[string]any)[vs.indexName]
	if !ok {
		return nil, fmt.Errorf("invalid response: no data for index %s", vs.indexName)
	}

	items, _ := data.([]any)
	docs := make([]schema.Document, len(items))

	for i, item := range items {
		metadata, _ := item.(map[string]any)

		docs[i] = schema.Document{
			PageContent: metadata[vs.opts.TextKey].(string),
			Metadata:    metadata,
		}
	}

	return docs, nil
}
