package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/documentloader"
	"github.com/hupe1980/golc/llm"
	"github.com/hupe1980/golc/textsplitter"
)

func main() {
	ctx := context.Background()

	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	llmSummarizationChain, err := chain.NewRefineSummarizationChain(openai)
	if err != nil {
		log.Fatal(err)
	}

	doc := `Large Language Models (LLMs) revolutionize natural language processing by providing 
	powerful tools for understanding, generating, and manipulating text at an unprecedented scale.
	
	Discover the limitless possibilities of Large Language Models (LLMs), advanced AI models 
	capable of understanding and generating human-like text across various domains and languages.

	Harness the power of state-of-the-art Large Language Models (LLMs) to enhance your applications 
	with advanced natural language processing capabilities, enabling tasks such as chatbots, 
	translation, sentiment analysis, and more.
	`

	loader := documentloader.NewTextLoader(strings.NewReader(doc))

	docs, err := loader.LoadAndSplit(ctx, textsplitter.NewRecusiveCharacterTextSplitter())
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llmSummarizationChain.Call(ctx, map[string]any{"inputDocuments": docs})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion["text"])
}
