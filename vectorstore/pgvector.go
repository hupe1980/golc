package vectorstore

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/hupe1980/golc/schema"
// 	pgv "github.com/pgvector/pgvector-go"
// )

// // Compile time check to ensure PGVector satisfies the VectorStore interface.
// var _ schema.VectorStore = (*PGVector)(nil)

// type PGVectorOptions struct {
// 	DriverName string
// }

// type PGVector struct {
// 	db       *sql.DB
// 	embedder schema.Embedder
// 	opts     PGVectorOptions
// }

// func NewPGVector(dataSourceName string, embedder schema.Embedder, optFns ...func(*PGVectorOptions)) (*PGVector, error) {
// 	opts := PGVectorOptions{
// 		DriverName: "pgx",
// 	}

// 	for _, fn := range optFns {
// 		fn(&opts)
// 	}

// 	db, err := sql.Open(opts.DriverName, dataSourceName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &PGVector{
// 		db:       db,
// 		embedder: embedder,
// 		opts:     opts,
// 	}, nil
// }

// func (vs *PGVector) AddDocuments(ctx context.Context, docs []schema.Document) error {
// 	texts := make([]string, len(docs))
// 	for i, doc := range docs {
// 		texts[i] = doc.PageContent
// 	}

// 	vectors, err := vs.embedder.EmbedDocuments(ctx, texts)
// 	if err != nil {
// 		return err
// 	}

// 	for i, doc := range docs {
// 		data := map[string]any{}
// 		data[store.embeddingStoreColumnName] = pgv.NewVector(float64ToFloat32(vectors[i]))
// 		data[store.textColumnName] = doc.PageContent

// 		if _, err := vs.db.ExecContext(ctx, ""); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (vs *PGVector) SimilaritySearch(ctx context.Context, query string) ([]schema.Document, error) {
// 	return nil, nil
// }

// func (vs *PGVector) CreateVectorExtensionIfNotExists(ctx context.Context) error {
// 	tx, err := vs.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	// Inspired by https://github.com/langchain-ai/langchain/blob/v0.0.340/libs/langchain/langchain/vectorstores/pgvector.py#L167
// 	// The advisor lock fixes issue arising from concurrent
// 	// creation of the vector extension.
// 	// https://github.com/langchain-ai/langchain/issues/12933
// 	// For more information see:
// 	// https://www.postgresql.org/docs/16/explicit-locking.html#ADVISORY-LOCKS
// 	if _, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(1573678846307946495)"); err != nil {
// 		return err
// 	}

// 	if _, err := tx.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
// 		return err
// 	}

// 	return tx.Commit()
// }

// func (vs *PGVector) Close() error {
// 	return vs.db.Close()
// }
