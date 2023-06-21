package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/agent"
	"github.com/hupe1980/golc/integration"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/tool"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	wikiTool := tool.NewWikipedia(integration.NewWikipedia())

	agent, err := agent.New(openai, []schema.Tool{wikiTool}, agent.ZeroShotReactDescriptionAgentType)
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(context.Background(), agent, "Who lived longer, Muhammad Ali or Alan Turing?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
