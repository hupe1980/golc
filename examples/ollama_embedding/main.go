package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hupe1980/golc/embedding"
	"github.com/hupe1980/golc/integration/ollama"
)

func main() {
	client := ollama.New("http://localhost:11434")

	embedder := embedding.NewOllama(client)

	e, err := embedder.EmbedText(context.Background(), "Hello llama2!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
