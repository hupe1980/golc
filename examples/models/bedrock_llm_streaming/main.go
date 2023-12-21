package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	client := bedrockruntime.NewFromConfig(cfg)

	bedrock, err := llm.NewBedrockAmazon(client, func(o *llm.BedrockAmazonOptions) {
		o.Callbacks = []schema.Callback{callback.NewStreamWriterHandler()}
		o.Stream = true
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := model.GeneratePrompt(context.Background(), bedrock, prompt.StringPromptValue("Write me a song about sparkling water.")); err != nil {
		log.Fatal(err)
	}
}
