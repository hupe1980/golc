// Package chatmodel provides a framework for working with chat-based large language models (LLMs).
package chatmodel

import "github.com/hupe1980/golc/schema"

func newChatGeneraton(text string, extFns ...func(o *schema.ChatMessageExtension)) schema.Generation { // nolint uparam
	return schema.Generation{
		Text:    text,
		Message: schema.NewAIChatMessage(text, extFns...),
	}
}
