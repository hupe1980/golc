package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	mathChain, err := chain.NewLLMMathFromLLM(openai)
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(context.Background(), mathChain, "What is 13 raised to the .3432 power?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
