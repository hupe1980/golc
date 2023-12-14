---
title: Sagemaker Endpoint
description: All about Sagemaker Endpoint.
weight: 70
---

1. Create a ContentHandler for Input/Output Transformation 
```go
type ModelRequest struct{
    Input string `json:"input"`
}

type ModelResponse []struct{
    GeneratedText string `json:"generated_text"`
}

type Transformer struct{}

func (mt *Transformer) TransformInput(prompt string) ([]byte, error) {
    return json.Marshal(&ModelRequest{
        Input: prompt,
    })
}

func (mt *Transformer) TransformOutput(output []byte) (string, error) {
    var res ModelResponse
    err := json.Unmarshal(output, &res)
    if err != nil {
        return "", err
    }

    return res[0].GeneratedText, nil
}

contentHandler := NewContentHandler("application/json", "application/json", Transformer{})
```

2. Create the Sagemaker Endpoint LLM
```go
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
  // Error handling
}

client := sagemakerruntime.NewFromConfig(cfg)

endpoint, err := NewSagemakerEndpoint(client, "my-endpoint", contentHandler)
if err != nil {
  // Error handling
}
```