package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/hupe1980/golc/texttospeech"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client := polly.NewFromConfig(cfg)

	polly := texttospeech.NewAmazonPolly(client)

	stream, err := polly.SynthesizeSpeech(context.Background(), "Hello world! My name is Joanna")
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Close()

	if err := stream.Play(); err != nil {
		log.Fatal(err)
	}
}
