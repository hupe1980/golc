package vectorstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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

	// IndexName is the name of the index to store the vectors.
	IndexName string

	// AdditionalFields is a list of additional fields to retrieve during similarity search.
	AdditionalFields []string
}

// Weaviate represents a Weaviate vector store.
type Weaviate struct {
	client   *weaviate.Client
	embedder schema.Embedder
	opts     WeaviateOptions
}

// NewWeaviate creates a new Weaviate vector store with the given Weaviate client, embedder, and optional configuration options.
func NewWeaviate(client *weaviate.Client, embedder schema.Embedder, optFns ...func(*WeaviateOptions)) *Weaviate {
	opts := WeaviateOptions{
		TextKey:   "text",
		TopK:      4,
		IndexName: fmt.Sprintf("GoLC_%s", uuid.New().String()),
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Weaviate{
		client:   client,
		embedder: embedder,
		opts:     opts,
	}
}

// CreateClassIfNotExist checks if the Weaviate class for the vector store exists, and creates it if it doesn't.
func (vs *Weaviate) CreateClassIfNotExist(ctx context.Context) error {
	exist, err := vs.client.Schema().ClassExistenceChecker().WithClassName(vs.opts.IndexName).Do(ctx)
	if err != nil {
		return err
	}

	if !exist {
		if ccErr := vs.client.Schema().ClassCreator().WithClass(&models.Class{
			Class: vs.opts.IndexName,
			Properties: []*models.Property{
				{
					Name:     vs.opts.TextKey,
					DataType: []string{"text"},
				},
			},
		}).Do(ctx); ccErr != nil {
			return ccErr
		}
	}

	return nil
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
			Class:      vs.opts.IndexName,
			ID:         strfmt.UUID(uuid.New().String()),
			Vector:     util.Float64ToFloat32(vectors[i]),
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

	nearVector := vs.client.GraphQL().NearVectorArgBuilder().WithVector(util.Float64ToFloat32(vector))

	fields := []graphql.Field{
		{Name: vs.opts.TextKey},
	}

	for _, fieldName := range vs.opts.AdditionalFields {
		fields = append(fields, graphql.Field{
			Name: fieldName,
		})
	}

	res, err := vs.client.GraphQL().
		Get().
		WithNearVector(nearVector).
		WithClassName(vs.opts.IndexName).
		WithFields(fields...).
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

	data, ok := res.Data["Get"].(map[string]any)[vs.opts.IndexName]
	if !ok {
		return nil, fmt.Errorf("invalid response: no data for index %s", vs.opts.IndexName)
	}

	items, _ := data.([]any)
	docs := make([]schema.Document, len(items))

	for i, item := range items {
		metadata, _ := item.(map[string]any)

		docs[i] = schema.Document{
			PageContent: metadata[vs.opts.TextKey].(string),
		}

		for _, field := range vs.opts.AdditionalFields {
			if v, ok := metadata[field]; ok {
				docs[i].Metadata[field] = v
			}
		}
	}

	return docs, nil
}

// Delete removes a document from the Weaviate vector store based on its UUID.
func (vs *Weaviate) Delete(ctx context.Context, uuid string) error {
	return vs.client.Data().Deleter().WithID(uuid).Do(ctx)
}
