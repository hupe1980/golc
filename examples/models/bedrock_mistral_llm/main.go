package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/model"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/prompt"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	client := bedrockruntime.NewFromConfig(cfg)

	bedrock, err := llm.NewBedrockMistral(client)
	if err != nil {
		log.Fatal(err)
	}

	res, err := model.GeneratePrompt(context.Background(), bedrock, prompt.StringPromptValue("Tell me a joke"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
