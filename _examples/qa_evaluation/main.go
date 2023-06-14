package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/evaluation"
	"github.com/hupe1980/golc/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	qaEval, err := evaluation.NewQAEvalChain(openai)
	if err != nil {
		log.Fatal(err)
	}

	examples := []map[string]string{{
		qaEval.QuestionKey(): "What is a LLM?",
		qaEval.AnswerKey():   "A Large Language Model (LLM) is a powerful AI model capable of understanding and generating human-like text at a large scale.",
	}}

	predictions := []map[string]string{{
		qaEval.PredictionKey(): "An apple is a round fruit with a crisp and juicy texture, typically red, green, or yellow in color, and often consumed as a snack or used in various culinary preparations.",
	}}

	grade, err := qaEval.Evaluate(context.Background(), examples, predictions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(grade[0]["text"])
}
