package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/embedding"
)

func main() {
	embedder, err := embedding.NewCohere(os.Getenv("COHERE_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	e, err := embedder.EmbedText(context.Background(), "Hello llama2!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
