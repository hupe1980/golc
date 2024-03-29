package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hupe1980/golc/embedding"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	client := bedrockruntime.NewFromConfig(cfg)

	embedder := embedding.NewBedrockAmazon(client)

	e, err := embedder.EmbedText(context.Background(), "Hello bedrock!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
