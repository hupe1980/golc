---
title: Retrieval Augmented Generation (RAG)
description: Interacting with data sources.
weight: 50
---

Retrieval Augmented Generation (RAG) is a powerful technique in the field of natural language processing that combines the capabilities of retrieval models and generation models. It involves interacting with data sources to retrieve relevant information and using it to enhance the generation process. RAG enables tasks such as summarization of lengthy text, question-answering based on specific datasets, and more.

The following example demonstrates an end-to-end workflow of retrieval augmented generation using GoLC:

1. The user makes a request to the GoLC app.
2. The app issues a search query to the retriever based on the user request.
3. The retriever returns search results with excerpts of relevant documents from the ingested enterprise data.
4. The app sends the user request and along with the data retrieved from the index as context in the LLM prompt.
5. The LLM returns a succinct response to the user request based on the retrieved data.
6. The response from the LLM is sent back to the user.

For detailed usage instructions and examples of how to use the retrievers, see the following sections.