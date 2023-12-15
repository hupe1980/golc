// Package vectorstore provides functionality for storing and managing vector embeddings.
package vectorstore

import (
	"github.com/hupe1980/golc/retriever"
	"github.com/hupe1980/golc/schema"
)

// ToRetriever takes a vector store and returns a retriever
func ToRetriever(vectorStore schema.VectorStore, optFns ...func(o *retriever.VectorStoreOptions)) schema.Retriever {
	return retriever.NewVectorStore(vectorStore, optFns...)
}
