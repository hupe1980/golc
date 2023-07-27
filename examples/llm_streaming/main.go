package main

import (
	"context"
	"log"
	"os"

	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *llm.OpenAIOptions) {
		o.Stream = true
	})
	if err != nil {
		log.Fatal(err)
	}

	_, mErr := model.GeneratePrompt(context.Background(), openai, prompt.StringPromptValue("Write me a song about sparkling water."), func(o *model.Options) {
		o.Callbacks = []schema.Callback{callback.NewStreamingStdOutHandler()}
	})
	if mErr != nil {
		log.Fatal(mErr)
	}
}
