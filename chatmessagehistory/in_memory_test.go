package chatmessagehistory

import (
	"context"
	"testing"

	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestInMemory(t *testing.T) {
	t.Run("Messages", func(t *testing.T) {
		// Create a test InMemory instance
		inMemory := NewInMemory()

		t.Run("Messages returns empty history initially", func(t *testing.T) {
			// Call the Messages method
			messages, err := inMemory.Messages(context.TODO())

			// Assert that an empty history is returned
			assert.NoError(t, err)
			assert.Empty(t, messages)
		})

		t.Run("Messages returns added messages", func(t *testing.T) {
			// Add some messages
			message1 := schema.NewHumanChatMessage("Message 1")
			message2 := schema.NewAIChatMessage("Message 2")
			err := inMemory.AddMessage(context.TODO(), message1)
			assert.NoError(t, err)
			err = inMemory.AddMessage(context.TODO(), message2)
			assert.NoError(t, err)

			// Call the Messages method
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)

			// Assert that the expected messages are returned
			expectedMessages := []schema.ChatMessage{message1, message2}
			assert.ElementsMatch(t, expectedMessages, messages)
		})
	})

	t.Run("AddUserMessage", func(t *testing.T) {
		// Create a test InMemory instance
		inMemory := NewInMemory()

		t.Run("AddUserMessage adds a user message", func(t *testing.T) {
			// Call the AddUserMessage method
			err := inMemory.AddUserMessage(context.TODO(), "User message")
			assert.NoError(t, err)

			// Assert that the message was added successfully
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)
			assert.Len(t, messages, 1)
		})

		t.Run("AddUserMessage creates a human chat message", func(t *testing.T) {
			// Call the AddUserMessage method
			_ = inMemory.AddUserMessage(context.TODO(), "User message")

			// Get the last message
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)

			lastMessage := messages[len(messages)-1]

			// Assert that the last message is a human chat message
			_, isHumanChatMessage := lastMessage.(*schema.HumanChatMessage)
			assert.True(t, isHumanChatMessage)
		})
	})

	t.Run("AddAIMessage", func(t *testing.T) {
		// Create a test InMemory instance
		inMemory := NewInMemory()

		t.Run("AddAIMessage adds an AI message", func(t *testing.T) {
			// Call the AddAIMessage method
			err := inMemory.AddAIMessage(context.TODO(), "AI message")
			assert.NoError(t, err)
			// Assert that the message was added successfully
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)
			assert.Len(t, messages, 1)
		})

		t.Run("AddAIMessage creates an AI chat message", func(t *testing.T) {
			// Call the AddAIMessage method
			_ = inMemory.AddAIMessage(context.TODO(), "AI message")

			// Get the last message
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)

			lastMessage := messages[len(messages)-1]

			// Assert that the last message is an AI chat message
			_, isAIChatMessage := lastMessage.(*schema.AIChatMessage)
			assert.True(t, isAIChatMessage)
		})
	})

	t.Run("Clear", func(t *testing.T) {
		// Create a test InMemory instance
		inMemory := NewInMemoryWithMessages([]schema.ChatMessage{
			schema.NewHumanChatMessage("Message 1"),
			schema.NewAIChatMessage("Message 2"),
		})

		t.Run("Clear removes all messages", func(t *testing.T) {
			// Call the Clear method
			err := inMemory.Clear(context.TODO())
			assert.NoError(t, err)
			// Assert that all messages are cleared
			messages, err := inMemory.Messages(context.TODO())
			assert.NoError(t, err)
			assert.Empty(t, messages)
		})
	})
}
