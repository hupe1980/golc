package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hupe1980/golc/embedding"
)

func main() {
	embedder, err := embedding.NewCybertron(func(o *embedding.CybertronOptions) {
		o.Model = "sentence-transformers/all-MiniLM-L6-v2"
		o.ModelsDir = ".models"
	})
	if err != nil {
		log.Fatal(err)
	}

	e, err := embedder.EmbedText(context.Background(), "Hello cybertron!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e)
}
