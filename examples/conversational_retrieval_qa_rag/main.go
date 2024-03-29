package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/model/chatmodel"
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

func (r *mockRetriever) Verbose() bool {
	return false
}

func (r *mockRetriever) Callbacks() []schema.Callback {
	return nil
}

func main() {
	golc.Verbose = true

	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	conversationalRetrievalQAChain, err := rag.NewConversationalRetrievalQA(openai, &mockRetriever{}, func(o *rag.ConversationalRetrievalQAOptions) {
		o.ReturnGeneratedQuestion = true
	})
	if err != nil {
		log.Fatal(err)
	}

	question1 := "Why don't scientists trust atoms?"

	result1, err := golc.Call(context.Background(), conversationalRetrievalQAChain, schema.ChainValues{
		"question": question1,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[i] Question:", question1)
	fmt.Println("[i] Generated Question:", result1["generatedQuestion"])
	fmt.Println("[i] Answer:", result1["answer"])

	question2 := "Can you explain it better?"

	result2, err := golc.Call(context.Background(), conversationalRetrievalQAChain, schema.ChainValues{
		"question": question2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[i] Question:", question2)
	fmt.Println("[i] Generated Question:", result2["generatedQuestion"])
	fmt.Println("[i] Answer:", result2["answer"])
}
