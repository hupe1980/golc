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
	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Temperature = 0.2 // optional
	})
	if err != nil {
		log.Fatal(err)
	}

	pv := prompt.StringPromptValue("What year was Einstein born?")

	result, err := model.GeneratePrompt(context.Background(), openai, pv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Generations[0].Text) // Output: Einstein was born in 1879.
}
