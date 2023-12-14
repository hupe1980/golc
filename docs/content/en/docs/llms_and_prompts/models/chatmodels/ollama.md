---
title: Ollama
description: All about Ollama.
weight: 40
---

```go
client := ollama.New("http://localhost:11434")

model, err := chatmodel.NewOllama(client, func(o *llm.OllamaOptions) {
    o.ModelName = "llama2"
})
if err != nil {
    // Error handling
}
```