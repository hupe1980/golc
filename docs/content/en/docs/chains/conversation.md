---
title: Conversation
description: All about conversation chains.
weight: 50
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
	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	conversationChain, err := chain.NewConversation(openai)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	result1, err := golc.SimpleCall(ctx, conversationChain, "What year was Einstein born?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result1)

	result2, err := golc.SimpleCall(ctx, conversationChain, "Multiply the year by 3.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result2)
}
```
Output:
```text
Einstein was born in 1879.
1879 multiplied by 3 equals 5637.
```