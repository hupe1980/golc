# 🦜️🔗 GoLC
![Build Status](https://github.com/hupe1980/golc/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/golc.svg)](https://pkg.go.dev/github.com/hupe1980/golc)
[![goreportcard](https://goreportcard.com/badge/github.com/hupe1980/golc)](https://goreportcard.com/report/github.com/hupe1980/golc)
[![codecov](https://codecov.io/gh/hupe1980/golc/branch/main/graph/badge.svg?token=Y4N7H8557X)](https://codecov.io/gh/hupe1980/golc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

🚀 Building Go applications with LLMs through composability
> GoLC is an innovative project heavily inspired by the [LangChain](https://github.com/hwchase17/langchain/tree/master) project, aimed at building applications with Large Language Models (LLMs) by leveraging the concept of composability. It provides a framework that enables developers to create and integrate LLM-based applications seamlessly. Through the principles of composability, GoLC allows for the modular construction of LLM-based components, offering flexibility and extensibility to develop powerful language processing applications. By leveraging the capabilities of LLMs and embracing composability, GoLC brings new opportunities to the Golang ecosystem for the development of natural language processing applications.

## Features
GoLC offers a range of features to enhance the development of language processing applications:

- 📃 LLMs and Prompts: GoLC simplifies the management and optimization of prompts and provides a generic interface for working with Large Language Models (LLMs). This simplifies the utilization of LLMs in your applications.
- 🔗 Chains: GoLC enables the creation of sequences of calls to LLMs or other utilities. It provides a standardized interface for chains, allowing for seamless integration with various tools. Additionally, GoLC offers pre-built end-to-end chains designed for common application scenarios, saving development time and effort.
- 📚 Retrieval Augmented Generation (RAG): GoLC supports specific types of chains that interact with data sources. This functionality enables tasks such as summarization of lengthy text and question-answering based on specific datasets. With GoLC, you can leverage RAG capabilities to enhance your language processing applications.
- 🤖 Agents: GoLC empowers the creation of agents that leverage LLMs to make informed decisions, take actions, observe results, and iterate until completion. By incorporating agents into your applications, you can enhance their intelligence and adaptability.
- 🧠 Memory: GoLC includes memory functionality that facilitates the persistence of state between chain or agent calls. This feature allows your applications to maintain context and retain important information throughout the processing pipeline. GoLC provides a standardized memory interface along with a selection of memory implementations for flexibility.
- 🎓 Evaluation: GoLC simplifies the evaluation of generative models, which are traditionally challenging to assess using conventional metrics. By utilizing language models themselves for evaluation, GoLC provides a novel approach to assessing the performance of generative models.
- 🚓 Moderation: GoLC incorporates essential moderation functionalities to enhance the security and appropriateness of language processing applications. This includes prompt injection detection, detection and redaction of Personally Identifiable Information (PII), identification of toxic content, and more.
- 📄 Document Processing: GoLC provides comprehensive document processing capabilities, including loading, transforming, and compressing. It offers a versatile set of tools to streamline document-related tasks, making it an ideal solution for document-centric language processing applications.

## Installation
Use Go modules to include golc in your project:
```
go get github.com/hupe1980/golc
```

## Usage
```golang
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/chatmodel"
)

func main() {
	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
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

For more example usage, see [examples](./examples).

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or new features you would like to see.

## References
- https://github.com/langchain-ai/langchain/
- https://www.promptingguide.ai/

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.


