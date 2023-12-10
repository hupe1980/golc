---
title: Anthropic
description: All about Anthropic.
weight: 10
---

```go
anthropic, err := chatmodel.NewAnthropic(os.Getenv("ANTHROPIC_API_KEY"))
if err != nil {
   // Error handling
}
```