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

func float64ToFloat32(v []float64) []float32 {
	v32 := make([]float32, len(v))
	for i, f := range v {
		v32[i] = float32(f)
	}

	return v32
}
