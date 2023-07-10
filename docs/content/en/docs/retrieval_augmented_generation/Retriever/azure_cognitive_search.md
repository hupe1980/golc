---
title: Azure Cognitive Search
description: A cloud-based search service that empowers developers to create intelligent search experiences for various applications and data sources.
weight: 10
---

Azure Cognitive Search is a cloud-based search service provided by Microsoft that offers developers a comprehensive set of tools and APIs to build intelligent search experiences. It enables organizations to create powerful search solutions for various applications, including websites, e-commerce platforms, and internal knowledge bases.

```go
apiKey := "Your API Key"
serviceName := "The serive name"
indexName := "The index name"

// Create a retriever
r := retriever.NewAzureCognitiveSearch(apiKey, serviceName, indexName)

// Now you can use retrieved documents from the Azure Cognitive Search index
docs, err := r.GetRelevantDocuments(context.Backgound(), "Your query")
if err != nil {
    // Handle error
}
```