---
title: Google GenAI
description: All about Google GenAI.
weight: 40
---

```go
ctx := context.Background()

client, err := generativelanguage.NewGenerativeClient(ctx)
if err != nil {
    // Error handling
}

defer client.Close()

llm, err := llm.NewGoogleGenAI(client)
if err != nil {
   // Error handling
}
```