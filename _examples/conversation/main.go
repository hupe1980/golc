package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/llm"
	"github.com/hupe1980/golc/schema"
)

func main() {
	golc.Verbose = true

	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	conversationChain, err := chain.NewConversation(openai, func(o *chain.ConversationOptions) {
		o.Callbacks = []schema.Callback{callback.NewStdOutHandler()}
	})
	if err != nil {
		log.Fatal(err)
	}

	result1, err := conversationChain.Run(context.Background(), "What is the meaning of life?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result1)

	result2, err := conversationChain.Run(context.Background(), "What is 1+1?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result2)
}
