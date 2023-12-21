package main

import (
	"context"
	"log"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/integration/ollama"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

// Start ollama
// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
// docker exec -it ollama ollama run llama2

func main() {
	client := ollama.New("http://localhost:11434")

	llm, err := llm.NewOllama(client, func(o *llm.OllamaOptions) {
		o.ModelName = "llama2"
		o.Stream = true
		o.Callbacks = []schema.Callback{callback.NewStreamWriterHandler()}
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := model.GeneratePrompt(context.Background(), llm, prompt.StringPromptValue("Hello llama2!")); err != nil {
		log.Fatal(err)
	}
}
