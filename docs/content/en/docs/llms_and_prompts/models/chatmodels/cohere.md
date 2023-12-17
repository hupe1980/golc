---
title: Cohere
description: All about Cohere.
weight: 30
---

```go
cohere, err := chatmodel.NewCohere(os.Getenv("COHERE_API_KEY"))
if err != nil {
   // Error handling
}
```

## Streaming
```go
cohere, err := chatmodel.NewCohere(os.Getenv("COHERE_API_KEY"), func(o *chatmodel.CohereOptions) {
   o.Callbacks = []schema.Callback{callback.NewStreamWriterHandler()}
   o.Stream = true
})
if err != nil {
   // Error handling
}
```