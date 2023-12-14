---
title: Gemini
description: All about Gemini.
weight: 40
---

```go
ctx := context.Background()

client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
if err != nil {
    // Error handling
}

defer client.Close()

llm, err := llm.NewGemini(client)
if err != nil {
   // Error handling
}
```