package vectorstore

import (
	"github.com/hupe1980/golc/retriever"
	"github.com/hupe1980/golc/schema"
)

// ToRetriever takes a vector store and returns a retriever
func ToRetriever(vectorStore schema.VectorStore) schema.Retriever {
	return retriever.NewVectorStore(vectorStore)
}
