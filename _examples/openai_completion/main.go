package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	completion, err := openai.Call(context.Background(), "What is the capital of France?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
