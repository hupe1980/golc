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

type Person struct {
	Name    string `json:"name" description:"The person's name"`
	Age     int    `json:"age" description:"The person's age"`
	FavFood string `json:"fav_food,omitempty" description:"The person's favorite food"`
}

func main() {
	golc.Verbose = true

	chatModel, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.ModelName = "gpt-4"
		o.Temperature = 0
	})
	if err != nil {
		log.Fatal(err)
	}

	structuredOutputChain, err := chain.NewStructuredOutput(chatModel, []chain.OutputCandidate{
		{
			Name:        "Person",
			Description: "Information about a person",
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
