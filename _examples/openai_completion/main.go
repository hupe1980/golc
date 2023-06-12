package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/llm/openai"
)

func main() {
	llm, err := openai.New(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(context.Background(), "What is the capital of France?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
