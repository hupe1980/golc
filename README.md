# 🦜️🔗 GoLC

⚡ Building applications with LLMs through composability ⚡

![Build Status](https://github.com/hupe1980/golc/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/golc.svg)](https://pkg.go.dev/github.com/hupe1980/golc)
> GoLC is an innovative project heavily inspired by the [LangChain](https://github.com/hwchase17/langchain/tree/master) project, aimed at building applications with Large Language Models (LLMs) by leveraging the concept of composability. It provides a framework that enables developers to create and integrate LLM-based applications seamlessly. Through the principles of composability, GoLC allows for the modular construction of LLM-based components, offering flexibility and extensibility to develop powerful language processing applications. By leveraging the capabilities of LLMs and embracing composability, GoLC brings new opportunities to the Golang ecosystem for the development of natural language processing applications.

## How to use
```golang
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/llm/openai"
)

func main() {
	llm, err := openai.New(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(context.Background(), "What is the capital of France?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
```
Output:
```text
The capital of France is Paris.
```

For more example usage, see [_examples](./_examples).

## References
- https://github.com/hwchase17/langchain/tree/master
- https://www.promptingguide.ai/

## License
[MIT](LICENCE)

