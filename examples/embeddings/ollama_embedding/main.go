package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hupe1980/golc/embedding"
	"github.com/hupe1980/golc/integration/ollama"
)

// Start ollama
// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
// docker exec -it ollama ollama run llama2

func main() {
	client := ollama.New("http://localhost:11434")

	embedder := embedding.NewOllama(client)

	e, err := embedder.EmbedText(context.Background(), "Hello llama2!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
