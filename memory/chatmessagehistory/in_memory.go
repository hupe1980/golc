package chatmessagehistory

import "github.com/hupe1980/golc"

// Compile time check to ensure InMemory satisfies the ChatMessageHistory interface.
var _ golc.ChatMessageHistory = (*InMemory)(nil)

type InMemory struct {
	messages []golc.ChatMessage
}

func NewInMemory() *InMemory {
	return &InMemory{
		messages: []golc.ChatMessage{},
	}
}

func NewInMemoryWithMessages(messages []golc.ChatMessage) *InMemory {
	return &InMemory{
		messages: messages,
	}
}

func (mh *InMemory) Messages() ([]golc.ChatMessage, error) {
	return mh.messages, nil
}

func (mh *InMemory) AddUserMessage(text string) error {
	message := golc.NewHumanChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *InMemory) AddAIMessage(text string) error {
	message := golc.NewAIChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *InMemory) AddMessage(message golc.ChatMessage) error {
	mh.messages = append(mh.messages, message)
	return nil
}

func (mh *InMemory) Clear() error {
	mh.messages = []golc.ChatMessage{}
	return nil
}
