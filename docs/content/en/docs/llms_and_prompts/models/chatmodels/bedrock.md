---
title: Amazon Bedrock
description: All about Amazon Bedrock.
weight: 20
---

```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := chatmodel.NewBedrock(client)
if err != nil {
    // Error handling
}
```

## Anthrophic Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := chatmodel.NewBedrockAntrophic(client)
if err != nil {
    // Error handling
}
```

## Meta Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := chatmodel.NewBedrockMeta(client)
if err != nil {
    // Error handling
}
```