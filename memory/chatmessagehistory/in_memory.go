package chatmessagehistory

import "github.com/hupe1980/golc/schema"

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

func (mh *InMemory) Messages() (schema.ChatMessages, error) {
	return mh.messages, nil
}

func (mh *InMemory) AddUserMessage(text string) error {
	message := schema.NewHumanChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *InMemory) AddAIMessage(text string) error {
	message := schema.NewAIChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *InMemory) AddMessage(message schema.ChatMessage) error {
	mh.messages = append(mh.messages, message)
	return nil
}

func (mh *InMemory) Clear() error {
	mh.messages = []schema.ChatMessage{}
	return nil
}
