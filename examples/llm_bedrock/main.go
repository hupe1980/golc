package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/model/llm"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client := bedrockruntime.NewFromConfig(cfg)

	bedrock, err := llm.NewBedrockAntrophic(client)
	if err != nil {
		log.Fatal(err)
	}

	res, err := bedrock.Generate(context.Background(), "These are a few of my favorite")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Generations[0].Text)
}
