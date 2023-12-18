package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
)

func main() {
	client := ollama.New("http://localhost:11434")

	llm, err := llm.NewOllama(client, func(o *llm.OllamaOptions) {
		o.ModelName = "llama2"
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.GeneratePrompt(context.Background(), llm, prompt.StringPromptValue("Hello llama2!"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
