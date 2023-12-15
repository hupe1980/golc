package vectorstore

import (
	"context"
	"fmt"

	"github.com/hupe1980/golc/integration/pinecone"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Pinecone satisfies the VectorStore interface.
var _ schema.VectorStore = (*Pinecone)(nil)

type PineconeOptions struct {
	Namespace string
	TopK      int64
}

type Pinecone struct {
	client   pinecone.Client
	embedder schema.Embedder
	textKey  string
	opts     PineconeOptions
}

func NewPinecone(client pinecone.Client, embedder schema.Embedder, textKey string, optFns ...func(*PineconeOptions)) (*Pinecone, error) {
	opts := PineconeOptions{
		TopK: 4,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Pinecone{
		client:   client,
		embedder: embedder,
		textKey:  textKey,
		opts:     opts,
	}, nil
}

func (vs *Pinecone) AddDocuments(ctx context.Context, docs []schema.Document) error {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.PageContent
	}

	vectors, err := vs.embedder.BatchEmbedText(ctx, texts)
	if err != nil {
		return err
	}

	metadata := make([]map[string]any, 0, len(docs))

	for i := 0; i < len(docs); i++ {
		m := make(map[string]any, len(docs[i].Metadata))
		for key, value := range docs[i].Metadata {
			m[key] = value
		}

		m[vs.textKey] = texts[i]

		metadata = append(metadata, m)
	}

	pineconeVectors, err := pinecone.ToPineconeVectors(vectors, metadata)
	if err != nil {
		return err
	}

	req := &pinecone.UpsertRequest{
		Vectors: pineconeVectors,
	}

	if vs.opts.Namespace != "" {
		req.Namespace = vs.opts.Namespace
	}

	_, err = vs.client.Upsert(ctx, req)

	return err
}

func (vs *Pinecone) SimilaritySearch(ctx context.Context, query string) ([]schema.Document, error) {
	vector, err := vs.embedder.EmbedText(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err := vs.client.Query(ctx, &pinecone.QueryRequest{
		Namespace:       vs.opts.Namespace,
		TopK:            vs.opts.TopK,
		IncludeMetadata: true,
		Vector:          vector,
	})
	if err != nil {
		return nil, err
	}

	docs := make([]schema.Document, 0, len(res.Matches))

	for _, match := range res.Matches {
		pageContent, ok := match.Metadata[vs.textKey].(string)
		if !ok {
			return nil, fmt.Errorf("no content for textKey %s", vs.textKey)
		}

		delete(match.Metadata, vs.textKey)

		doc := schema.Document{
			PageContent: pageContent,
			Metadata:    match.Metadata,
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
