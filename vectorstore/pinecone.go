package vectorstore

import (
	"context"

	"github.com/hupe1980/golc/integration/pinecone"
	"github.com/hupe1980/golc/schema"
	pc "github.com/pinecone-io/go-pinecone/pinecone_grpc"
)

// Compile time check to ensure Pinecone satisfies the VectorStore interface.
var _ schema.VectorStore = (*Pinecone)(nil)

type PineconeOptions struct {
	Namespace string
	UseGRPC   bool
}

type Pinecone struct {
	client   *pinecone.GRPCClient
	embedder schema.Embedder
	textKey  string
	opts     PineconeOptions
}

func NewPinecone(apiKey string, endpoint pinecone.Endpoint, embedder schema.Embedder, textKey string, optFns ...func(*PineconeOptions)) (*Pinecone, error) {
	opts := PineconeOptions{
		UseGRPC: false,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	client, err := pinecone.NewGRPCClient(apiKey, endpoint)
	if err != nil {
		return nil, err
	}

	return &Pinecone{
		client:   client,
		embedder: embedder,
		textKey:  textKey,
		opts:     opts,
	}, nil
}

func (vs *Pinecone) AddDocuments(ctx context.Context, docs []schema.Document) error {
	texts := make([]string, 0, len(docs))
	for _, doc := range docs {
		texts = append(texts, doc.PageContent)
	}

	vectors, err := vs.embedder.EmbedDocuments(ctx, texts)
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

	pineconeVectors, err := pinecone.ToPineconeGRPCVectors(vectors, metadata)
	if err != nil {
		return err
	}

	req := &pc.UpsertRequest{
		Vectors: pineconeVectors,
	}

	if vs.opts.Namespace != "" {
		req.Namespace = vs.opts.Namespace
	}

	_, err = vs.client.Upsert(ctx, req)

	return err
}

func (vs *Pinecone) SimilaritySearch(ctx context.Context, query string) ([]schema.Document, error) {
	return nil, nil
}
