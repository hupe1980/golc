---
title: Ollama
description: All about Ollama.
weight: 50
---

```go
client := ollama.New("http://localhost:11434")

llm, err := llm.NewOllama(client, func(o *llm.OllamaOptions) {
    o.ModelName = "llama2"
})
if err != nil {
    // Error handling
}
```