---
title: 🦜️🔗 GoLC
---

{{< blocks/cover title="🚀 Building Go applications with LLMs through composability" image_anchor="top" height="full" color="primary">}}
<a class="btn btn-lg btn-secondary me-3 mb-4" href="/golc/docs/">
  Learn More About GoLC <i class="fas fa-arrow-alt-circle-right ms-2"></i>
</a>
<p class="lead mt-5">GoLC is an innovative project heavily inspired by the <a class="link-dotted-underline" href="https://github.com/langchain-ai/langchain/">LangChain</a> project, aimed at building applications with Large Language Models (LLMs) by leveraging the concept of composability. It provides a framework that enables developers to create and integrate LLM-based applications seamlessly. Through the principles of composability, GoLC allows for the modular construction of LLM-based components, offering flexibility and extensibility to develop powerful language processing applications. By leveraging the capabilities of LLMs and embracing composability, GoLC brings new opportunities to the Golang ecosystem for the development of natural language processing applications.</p>
{{< blocks/link-down color="info" >}}
{{< /blocks/cover >}}


{{% blocks/lead color="light" %}}
GoLC offers a range of features to enhance the development of language processing applications:
{{% /blocks/lead %}}


{{% blocks/section color="primary" type="row" %}}

{{% blocks/feature icon="fas fa-terminal" title="LLMs and Prompts" %}}
GoLC simplifies the management and optimization of prompts and provides a generic interface for working with Large Language Models (LLMs). This simplifies the utilization of LLMs in your applications.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-link" title="Chains"  %}}
GoLC enables the creation of sequences of calls to LLMs or other utilities. It provides a standardized interface for chains, allowing for seamless integration with various tools. Additionally, GoLC offers pre-built end-to-end chains designed for common application scenarios, saving development time and effort.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-book" title="Retrieval Augmented Generation (RAG)" %}}
GoLC supports specific types of chains that interact with data sources. This functionality enables tasks such as summarization of lengthy text and question-answering based on specific datasets. With GoLC, you can leverage RAG capabilities to enhance your language processing applications.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-robot" title="Agents" %}}
GoLC empowers the creation of agents that leverage LLMs to make informed decisions, take actions, observe results, and iterate until completion. By incorporating agents into your applications, you can enhance their intelligence and adaptability.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-memory" title="Memory" %}}
GoLC includes memory functionality that facilitates the persistence of state between chain or agent calls. This feature allows your applications to maintain context and retain important information throughout the processing pipeline. GoLC provides a standardized memory interface along with a selection of memory implementations for flexibility.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-graduation-cap" title="Evaluation" %}}
GoLC simplifies the evaluation of generative models, which are traditionally challenging to assess using conventional metrics. By utilizing language models themselves for evaluation, GoLC provides a novel approach to assessing the performance of generative models.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-shield" title="Moderation" %}}
GoLC incorporates essential moderation functionalities to enhance the security and appropriateness of language processing applications. This includes prompt injection detection, detection and redaction of Personally Identifiable Information (PII), identification of toxic content, and more.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-file" title="Document Processing" %}}
GoLC provides comprehensive document processing capabilities, including loading, transforming, and compressing. It offers a versatile set of tools to streamline document-related tasks, making it an ideal solution for document-centric language processing applications.
{{% /blocks/feature %}}

{{% /blocks/section %}}