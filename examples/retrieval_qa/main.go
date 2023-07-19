package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/rag"
	"github.com/hupe1980/golc/schema"
)

type mockRetriever struct{}

func (r *mockRetriever) GetRelevantDocuments(ctx context.Context, query string) ([]schema.Document, error) {
	return []schema.Document{
		{PageContent: "Why don't scientists trust atoms? Because they make up everything!"},
		{PageContent: "Why did the bicycle fall over? Because it was two-tired!"},
	}, nil
}

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	retrievalQAChain, err := rag.NewRetrievalQA(openai, &mockRetriever{})
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(context.Background(), retrievalQAChain, "Why don't scientists trust atoms?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
