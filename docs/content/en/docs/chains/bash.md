---
title: Bash
description: All about bash chains.
weight: 60
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
	"github.com/hupe1980/golc/model/llm"
)

func main() {
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *llm.OpenAIOptions) {
		o.Temperature = 0.01
	})
	if err != nil {
		log.Fatal(err)
	}

	bashChain, err := chain.NewBash(openai)
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(context.Background(), bashChain, "Please write a bash script that prints 'Hello World' to the console.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
```
Output:
```text
Hello World
```