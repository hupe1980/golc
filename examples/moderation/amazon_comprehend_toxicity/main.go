package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/moderation"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	client := comprehend.NewFromConfig(cfg)

	moderationChain := moderation.NewAmazonComprehendToxicity(client)

	result, err := golc.SimpleCall(context.Background(), moderationChain, "I will kill you")
	if err != nil {
		log.Fatal(err) // toxic content found
	}

	fmt.Println(result)
}
