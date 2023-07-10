---
title: Getting Started
description: How to get started with GoLC.
weight: 20
---

## Installation
Use Go modules to include golc in your project:
```shell
go get github.com/hupe1980/golc
```

## Getting Predictions from Large Language Models
The core functionality of the GoLC project revolves around Language Models (LLMs), which excel at generating text based on input text. GoLC offers extensive support for a variety of pre-trained LLMs, providing developers with a wide range of options to choose from.

To leverage the power of LLMs in your application, you can initialize an LLM model, such as the OpenAI model, and make predictions. For example, you can use the OpenAI model to determine the birth year of Albert Einstein:
```go
import (
    "context"
    "os"

	"github.com/hupe1980/golc/model/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		// Handle error
	}

	prompt := "What year was Einstein born?"

	result, err := openai.Generate(context.Background(), prompt)
	if err != nil {
		// Handle error
	}

	fmt.Println(result.Generations[0].Text) // Output: Einstein was born in 1879.
}
```

GoLC also supports a wide range of other large language models out of the box. Additionally, it offers chat models that have been fine-tuned or specially trained to handle conversational interactions, such as chat messages or dialogue-based interactions. This expanded model support allows developers to build sophisticated applications with enhanced language processing capabilities.