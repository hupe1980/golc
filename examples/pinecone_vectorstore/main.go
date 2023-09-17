package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/embedding"
	"github.com/hupe1980/golc/integration/pinecone"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/vectorstore"
)

func main() {
	openai, err := embedding.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	client, err := pinecone.NewRestClient(os.Getenv("PINECONE_API_KEY"), pinecone.Endpoint{
		IndexName:   "golc",
		ProjectName: os.Getenv("PINECONE_PROJECT"),
		Environment: "us-west1-gcp-free",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	vs, err := vectorstore.NewPinecone(client, openai, "textKey", func(po *vectorstore.PineconeOptions) {
		po.TopK = 1
	})
	if err != nil {
		log.Fatal(err)
	}

	if err = vs.AddDocuments(context.Background(), []schema.Document{
		{
			PageContent: "Pizza is an Italian dish consisting of a flat, round base of dough topped with various ingredients, including tomato sauce, cheese, and various toppings.",
			Metadata: map[string]any{
				"cousine": "Italian",
			},
		},
		{
			PageContent: "Sushi is a Japanese dish consisting of vinegared rice combined with various ingredients, such as seafood, vegetables, and occasionally tropical fruits.",
			Metadata: map[string]any{
				"cousine": "Japanese",
			},
		},
		{
			PageContent: "A burger is a sandwich consisting of a cooked ground meat patty, usually beef, placed in a sliced bread roll or bun.",
			Metadata: map[string]any{
				"cousine": "American",
			},
		},
		{
			PageContent: "Pad Thai is a popular Thai stir-fried noodle dish made with rice noodles, eggs, tofu or shrimp, bean sprouts, peanuts, and lime.",
			Metadata: map[string]any{
				"cousine": "Thai",
			},
		},
		{
			PageContent: "Sashimi is a Japanese delicacy consisting of thinly sliced raw seafood, served with soy sauce and wasabi.",
			Metadata: map[string]any{
				"cousine": "Japanese",
			},
		},
		{
			PageContent: "Lasagna is an Italian pasta dish made with layers of pasta sheets, meat sauce, cheese, and bechamel sauce.",
			Metadata: map[string]any{
				"cousine": "Italian",
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	docs, err := vs.SimilaritySearch(context.Background(), "Cheeseburger")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(docs)
}
