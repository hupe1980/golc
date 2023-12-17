---
title: OpenAI
description: All about OpenAI.
weight: 70
---

```go
openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
if err != nil {
   // Error handling
}
```

## Streaming
```go
openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
   o.MaxTokens = 256
   o.Stream = true
   o.Callbacks = []schema.Callback{callback.NewStreamWriterHandler()}
})
if err != nil {
   // Error handling
}
```