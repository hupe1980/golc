---
title: Sagemaker Endpoint
description: All about Sagemaker Endpoint.
weight: 10
---

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

contentHandler := NewLLMContentHandler("application/json", "application/json", Transformer{})

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