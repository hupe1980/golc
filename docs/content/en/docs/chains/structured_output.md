---
title: StructuredOutput
description: All about structured output chains.
weight: 100
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
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

type Person struct {
	Name    string `json:"name" description:"The person's name"`
	Age     int    `json:"age" description:"The person's age"`
	FavFood string `json:"fav_food,omitempty" description:"The person's favorite food"`
}

func main() {
	chatModel, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Temperature = 0
	})
	if err != nil {
		log.Fatal(err)
	}

	pt := prompt.NewChatTemplate([]prompt.MessageTemplate{
		prompt.NewSystemMessageTemplate("You are a world class algorithm for extracting information in structured formats."),
		prompt.NewHumanMessageTemplate("Use the given format to extract information from the following input:\n{{.input}}\nTips: Make sure to answer in the correct format"),
	})

	structuredOutputChain, err := chain.NewStructuredOutput(chatModel, pt, []chain.OutputCandidate{
		{
			Name:        "Person",
			Description: "Identifying information about a person",
			Data:        &Person{},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.Call(context.Background(), structuredOutputChain, schema.ChainValues{
		"input": "Max is 21",
	})
	if err != nil {
		log.Fatal(err)
	}

	p, ok := result["output"].(*Person)
	if !ok {
		log.Fatal("output is not a person")
	}

	fmt.Println("Name:", p.Name)
	fmt.Println("Age:", p.Age)
}
```
Output:
```text
Name: Max
Age: 21
```