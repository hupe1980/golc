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

	bedrock, err := llm.NewBedrock(client, "amazon.titan-text-lite-v1", func(o *llm.BedrockOptions) {
		o.ModelParams = map[string]any{ // optional
			"temperature": 0.3,
		}
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
