package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/chatmodel/openai"
	"github.com/hupe1980/golc/schema"
)

func main() {
	llm, err := openai.New(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(context.Background(), []schema.Message{
		schema.SystemMessage{Text: "Hello, I am a friendly chatbot. I love to talk about movies, books and music. Answer in markdown format."},
		schema.HumanMessage{Text: "What would be a good company name for a company that makes colorful socks?"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion.GetText())
}
