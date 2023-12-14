---
title: Palm
description: All about Palm.
weight: 70
---

```go
ctx := context.Background()

client, err := generativelanguage.NewTextClient(ctx)
if err != nil {
    // Error handling
}

defer client.Close()

palm, err := llm.NewPalm(client)
if err != nil {
   // Error handling
}
```