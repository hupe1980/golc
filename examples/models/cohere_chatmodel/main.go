package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
)

func main() {
	cohere, err := chatmodel.NewCohere(os.Getenv("COHERE_API_KEY"), func(o *chatmodel.CohereOptions) {
		o.Temperature = 0.7 // optional
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.GeneratePrompt(context.Background(), cohere, prompt.StringPromptValue("How much cost the fish? A short answer please."))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
