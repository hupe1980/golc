---
title: Amazon Kendra
description: A powerful enterprise search service that allows you to easily find relevant information from your data using natural language queries.
weight: 10
---

Amazon Kendra is an intelligent search service provided by Amazon Web Services (AWS). It utilizes advanced natural language processing (NLP) and machine learning algorithms to enable powerful search capabilities across various data sources within an organization. Kendra is designed to help users find the information they need quickly and accurately, improving productivity and decision-making.

With Kendra, users can search across a wide range of content types, including documents, FAQs, knowledge bases, manuals, and websites. It supports multiple languages and can understand complex queries, synonyms, and contextual meanings to provide highly relevant search results.

```go
// Using the SDK's default configuration
cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
if err != nil {
    // Handle error
}

// Using the Config value, create the Kendra client
kendraClient := kendra.NewFromConfig(cfg)

// Using the client and index ID, create a retriever
r := retriever.NewAmazonKendra(kendraClient, "The Kendra Index ID")

// Now you can use retrieved documents from the Kendra index
docs, err := r.GetRelevantDocuments(context.Backgound(), "Your query")
if err != nil {
    // Handle error
}
```