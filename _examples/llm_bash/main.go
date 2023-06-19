package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	bashChain, err := chain.NewLLMBashFromLLM(openai)
	if err != nil {
		log.Fatal(err)
	}

	result, err := chain.SimpleCall(context.Background(), bashChain, "Please write a bash script that prints 'Hello World' to the console.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
