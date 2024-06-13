package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	client := bedrockruntime.NewFromConfig(cfg)

	bedrock, err := chatmodel.NewBedrock(client, "anthropic.claude-3-sonnet-20240229-v1:0", func(o *chatmodel.BedrockOptions) {
		o.Temperature = aws.Float32(0.3)
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.GeneratePrompt(context.Background(), bedrock, prompt.StringPromptValue("Hello ai!"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
