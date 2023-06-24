# 🦜️🔗 GoLC

⚡ Building applications with LLMs through composability in golang  ⚡

![Build Status](https://github.com/hupe1980/golc/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/golc.svg)](https://pkg.go.dev/github.com/hupe1980/golc)
> GoLC is an innovative project heavily inspired by the [LangChain](https://github.com/hwchase17/langchain/tree/master) project, aimed at building applications with Large Language Models (LLMs) by leveraging the concept of composability. It provides a framework that enables developers to create and integrate LLM-based applications seamlessly. Through the principles of composability, GoLC allows for the modular construction of LLM-based components, offering flexibility and extensibility to develop powerful language processing applications. By leveraging the capabilities of LLMs and embracing composability, GoLC brings new opportunities to the Golang ecosystem for the development of natural language processing applications.

## Installation
```
go get github.com/hupe1980/golc
```

## How to use
```golang
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/llm"
	"github.com/hupe1980/golc/schema"
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

For more example usage, see [_examples](./_examples).

## References
- https://github.com/hwchase17/langchain/tree/master
- https://www.promptingguide.ai/

## License
[MIT](LICENCE)

