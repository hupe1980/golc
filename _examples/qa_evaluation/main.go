package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/evaluation"
	"github.com/hupe1980/golc/llm/openai"
)

func main() {
	llm, err := openai.New(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	chain, err := evaluation.NewQAEvalChain(llm)
	if err != nil {
		log.Fatal(err)
	}

	examples := []map[string]string{{
		chain.QuestionKey(): "What is an LLM?",
		chain.AnswerKey():   "A Large Language Model (LLM) is a powerful AI model capable of understanding and generating human-like text at a large scale.",
	}}

	predictions := []map[string]string{{
		chain.PredictionKey(): "An apple is a round fruit with a crisp and juicy texture, typically red, green, or yellow in color, and often consumed as a snack or used in various culinary preparations.",
	}}

	grade, err := chain.Evaluate(context.Background(), examples, predictions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(grade[0]["text"])
}
