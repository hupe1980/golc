---
title: Amazon Bedrock
description: All about Amazon Bedrock.
weight: 20
---

{{< ghcode src="https://raw.githubusercontent.com/hupe1980/golc/main/examples/models/bedrock_llm/main.go" >}}

## Streaming
{{< ghcode src="https://raw.githubusercontent.com/hupe1980/golc/main/examples/models/bedrock_llm_streaming/main.go" >}}


## A121 Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := llm.NewBedrockA121(client)
if err != nil {
    // Error handling
}
```

## Amazon Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := llm.NewBedrockAmazon(client)
if err != nil {
    // Error handling
}
```

## Cohere Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := llm.NewBedrockCohere(client)
if err != nil {
    // Error handling
}
```

## Anthrophic Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := llm.NewBedrockAntrophic(client)
if err != nil {
    // Error handling
}
```

## Meta Support
```go
cfg, _ := config.LoadDefaultConfig(context.Background())
client := bedrockruntime.NewFromConfig(cfg)

bedrock, err := llm.NewBedrockMeta(client)
if err != nil {
    // Error handling
}
```