package chatmessagehistory

import (
	"context"

	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure InMemory satisfies the ChatMessageHistory interface.
var _ schema.ChatMessageHistory = (*InMemory)(nil)

type InMemory struct {
	messages schema.ChatMessages
}

func NewInMemory() *InMemory {
	return &InMemory{
		messages: schema.ChatMessages{},
	}
}

func NewInMemoryWithMessages(messages schema.ChatMessages) *InMemory {
	return &InMemory{
		messages: messages,
	}
}

func (mh *InMemory) Messages(ctx context.Context) (schema.ChatMessages, error) {
	return mh.messages, nil
}

func (mh *InMemory) AddUserMessage(ctx context.Context, text string) error {
	message := schema.NewHumanChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *InMemory) AddAIMessage(ctx context.Context, text string) error {
	message := schema.NewAIChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *InMemory) AddMessage(ctx context.Context, message schema.ChatMessage) error {
	mh.messages = append(mh.messages, message)
	return nil
}

func (mh *InMemory) Clear(ctx context.Context) error {
	mh.messages = []schema.ChatMessage{}
	return nil
}
