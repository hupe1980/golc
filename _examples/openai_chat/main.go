package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chatmodel"
)

func main() {
	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	completion, err := openai.Call(context.Background(), []golc.ChatMessage{
		golc.NewSystemChatMessage("Hello, I am a friendly chatbot. I love to talk about movies, books and music. Answer in markdown format."),
		golc.NewHumanChatMessage("What would be a good company name for a company that makes colorful socks?"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion.Text())
}
