---
title: Baidu Ernie
description: All about Baidu Ernie.
weight: 30
---

```go
ernie, err := chatmodel.NewErnie(os.Getenv("ERNIE_CLIENT_ID"), os.Getenv("ERNIE_CLIENT_SECRET"))
if err != nil {
   // Error handling
}
```