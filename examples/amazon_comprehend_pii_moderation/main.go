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

	moderationChain := moderation.NewAmazonComprehendPII(client)

	result, err := golc.SimpleCall(context.Background(), moderationChain, "My Name is Alfred E. Neuman")
	if err != nil {
		log.Fatal(err) // pii content found
	}

	fmt.Println(result)
}
