package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/documentloader"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/rag"
	"github.com/hupe1980/golc/textsplitter"
)

func main() {
	ctx := context.Background()

	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	llmSummarizationChain, err := rag.NewMapReduceSummarization(openai)
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

	loader := documentloader.NewText(strings.NewReader(doc))

	docs, err := loader.LoadAndSplit(ctx, textsplitter.NewRecusiveCharacterTextSplitter())
	if err != nil {
		log.Fatal(err)
	}

	completion, err := golc.SimpleCall(ctx, llmSummarizationChain, docs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
