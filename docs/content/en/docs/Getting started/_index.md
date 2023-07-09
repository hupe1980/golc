---
title: Getting Started
description: How to get started with GoLC
weight: 20
---

## Installation
Use Go modules to include golc in your project:
```shell
go get github.com/hupe1980/golc
```

## LLMs: Getting Predictions from Language Models
The core functionality of the GoLC project revolves around Language Models (LLMs), which excel at generating text based on input text. The following example shows the call of the openai model to determine the year of birth of Albert Einstein:
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