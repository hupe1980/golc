package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/model/llm"
)

func main() {
	model, err := llm.NewCohere(os.Getenv("COHERE_API_KEY"), func(o *llm.CohereOptions) {
		o.MaxTokens = 256
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.Generate(context.Background(), "How much cost the fish? A short answer please.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
