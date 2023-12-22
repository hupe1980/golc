package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/hupe1980/golc/documenttransformer"
	"github.com/hupe1980/golc/schema"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client := translate.NewFromConfig(cfg)

	doc := schema.Document{
		PageContent: "Das Pferd frisst keinen Gurkensalat",
	}

	at := documenttransformer.NewAmazonTranslate(client, "en", func(o *documenttransformer.AmazonTranslateOptions) {
		o.IncludeSourceText = true // optional
	})

	docs, err := at.Transform(context.Background(), []schema.Document{doc})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(docs[0].PageContent)
	fmt.Println(docs[0].Metadata)
}
