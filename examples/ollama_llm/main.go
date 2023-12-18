package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/model/llm"
)

func main() {
	client := ollama.New("http://localhost:11434")

	llm, err := llm.NewOllama(client, func(o *llm.OllamaOptions) {
		o.ModelName = "llama2"
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := llm.Generate(context.Background(), "Hello llama2!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
