---
title: Tagging
description: All about tagging chains.
weight: 110
---
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
)

type Tagging struct {
	Sentiment      string `json:"sentiment" enum:"'happy','neutral','sad'"`
	Aggressiveness int    `json:"aggressiveness" description:"describes how aggressive the statement is, the higher the number the more aggressive" enum:"1,2,3,4,5"`
	Language       string `json:"language" enum:"'spanish','english','french','german','italian'"`
}

func main() {
	chatModel, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Temperature = 0
	})
	if err != nil {
		log.Fatal(err)
	}

	taggingChain, err := chain.NewTagging(chatModel, &Tagging{})
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.Call(context.Background(), taggingChain, schema.ChainValues{
		"input": "Weather is ok here, I can go outside without much more than a coat",
	})
	if err != nil {
		log.Fatal(err)
	}

	t, ok := result["output"].(*Tagging)
	if !ok {
		log.Fatal("output is not tagging")
	}

	fmt.Println("Sentiment:", t.Sentiment)
	fmt.Println("Aggressiveness:", t.Aggressiveness)
	fmt.Println("Language:", t.Language)
}
```
Output:
```text
Sentiment: neutral
Aggressiveness: 3
Language: english
```