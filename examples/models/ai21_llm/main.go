package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/model/llm"
)

func main() {
	ai21, err := llm.NewAI21(os.Getenv("AI21_API_KEY"), func(o *llm.AI21Options) {
		o.MaxTokens = 256 // optional
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := ai21.Generate(context.Background(), "These are a few of my favorite")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
