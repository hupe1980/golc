package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/embedding"
)

func main() {
	embedder := embedding.NewCohere(os.Getenv("COHERE_API_KEY"))

	e, err := embedder.EmbedText(context.Background(), "Hello llama2!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
