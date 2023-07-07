package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/schema"
)

func main() {
	moderationChain, err := chain.NewOpenAIModeration(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.Call(context.Background(), moderationChain, schema.ChainValues{
		"input": "I will kill you",
	})
	if err != nil {
		log.Fatal(err) // content policy violation
	}

	fmt.Println(result)
}
