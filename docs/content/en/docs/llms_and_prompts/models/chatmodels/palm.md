---
title: Palm
description: All about Palm.
weight: 60
---

```go
ctx := context.Background()

// see https://pkg.go.dev/cloud.google.com/go/ai@v0.1.1/generativelanguage/apiv1beta2
c, err := generativelanguage.NewDiscussClient(ctx)
if err != nil {
	// Error handling
}
defer c.Close()

palm, err := chatmodel.NewPalm(c)
if err != nil {
   // Error handling
}
```