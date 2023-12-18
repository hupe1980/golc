package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
)

func main() {
	cohere, err := llm.NewCohere(os.Getenv("COHERE_API_KEY"), func(o *llm.CohereOptions) {
		o.MaxTokens = 256
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.GeneratePrompt(context.Background(), cohere, prompt.StringPromptValue("How much cost the fish? A short answer please."))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
