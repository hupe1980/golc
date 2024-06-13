package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	client := bedrockruntime.NewFromConfig(cfg)

	bedrock, err := chatmodel.NewBedrock(client, "amazon.titan-text-express-v1", func(o *chatmodel.BedrockOptions) {
		o.Temperature = aws.Float32(0.8)
		o.Callbacks = []schema.Callback{callback.NewStreamWriterHandler()}
		o.MaxTokens = aws.Int32(4096)
		o.Stream = true
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = model.GeneratePrompt(context.Background(), bedrock, prompt.StringPromptValue("Write me a song about sparkling water."))
	if err != nil {
		log.Fatal(err)
	}
}
