package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

func main() {
	info := callback.NewOpenAIHandler()

	openAI, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Callbacks = []schema.Callback{info}
	})
	if err != nil {
		log.Fatal(err)
	}

	t := prompt.NewSystemMessageTemplate("Hello World")

	pv, err := t.FormatPrompt(nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := model.GeneratePrompt(context.Background(), openAI, pv)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Generations[0].Text)

	fmt.Println(info)
}
