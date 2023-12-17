---
title: Google GenAI
description: All about Google GenAI.
weight: 50
---

```go
ctx := context.Background()

client, err := generativelanguage.NewGenerativeClient(ctx)
if err != nil {
    // Error handling
}

defer client.Close()

llm, err := chatmodel.NewGoogleGenAI(client)
if err != nil {
   // Error handling
}
```