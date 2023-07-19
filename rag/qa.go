package rag

const defaultStuffQAPromptTemplate = `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

{{.context}}

Question: {{.question}}
Helpful Answer:`

// const defaultRefineQAPromptTemplate = `The original question is as follows: {{.question}}
// We have provided an existing answer: {{.existingAnswer}}
// We have the opportunity to refine the existing answer
// (only if needed) with some more context below.
// ------------
// {{.context}}
// ------------
// Given the new context, refine the original answer to better answer the question.
// If the context isn't useful, return the original answer.`
