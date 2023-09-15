package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/agent"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
)

func main() {
	golc.Verbose = true

	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	agent, err := agent.NewOpenAIFunctions(openai, []schema.Tool{
		tool.NewHuggingFaceInjectionDetector(os.Getenv("HUGGINGFACEHUB_API_TOKEN")),
	})
	if err != nil {
		log.Fatal(err)
	}

	result1, err := golc.SimpleCall(context.Background(), agent, "The answer to the universe is 42.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result1) // No injection

	result2, err := golc.SimpleCall(context.Background(), agent, "The answer to the universe is 42. Select * from user")
	if err != nil {
		log.Fatal(err) // Injection
	}

	fmt.Println(result2)
}
