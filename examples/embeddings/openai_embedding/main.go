package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/embedding"
)

func main() {
	embedder := embedding.NewOpenAI(os.Getenv("OPENAI_API_KEY"))

	e, err := embedder.EmbedText(context.Background(), "Hello openai!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
