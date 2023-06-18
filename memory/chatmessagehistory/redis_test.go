package chatmessagehistory

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRedisClient struct {
	mock.Mock
}

func (c *mockRedisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	args := c.Called(ctx, key, start, stop)
	cmd := redis.NewStringSliceCmd(ctx)

	if val, ok := args.Get(0).([]string); ok {
		cmd.SetVal(val)
		return cmd
	}

	cmd.SetErr(redis.Nil)

	return cmd
}

func (c *mockRedisClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := c.Called(ctx, key, values[0])

	cmd := redis.NewIntCmd(ctx)

	if err, ok := args.Get(0).(error); ok {
		cmd.SetErr(err)
		return cmd
	}

	if i, ok := args.Get(0).(int64); ok {
		cmd.SetVal(i)
	}

	return cmd
}

func (c *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := c.Called(ctx, keys[0])

	cmd := redis.NewIntCmd(ctx)

	if err, ok := args.Get(0).(error); ok {
		cmd.SetErr(err)
		return cmd
	}

	return cmd
}

func (c *mockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := c.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func TestRedis(t *testing.T) {
	t.Run("Messages", func(t *testing.T) {
		mockClient := &mockRedisClient{}
		redisHistory := NewRedis(mockClient, "session1")

		t.Run("Messages returns chat messages", func(t *testing.T) {
			// Prepare Redis data
			messagesJSON := []string{
				`{"type":"human","text":"Message 1"}`,
				`{"type":"ai","text":"Message 2"}`,
			}
			mockClient.Mock = mock.Mock{}
			mockClient.On("LRange", mock.Anything, redisHistory.key(), int64(0), int64(-1)).
				Return(messagesJSON, nil)

			// Call the Messages method
			messages, err := redisHistory.Messages()
			assert.NoError(t, err)

			// Assert that the correct chat messages are returned
			expectedMessages := schema.ChatMessages{
				schema.NewHumanChatMessage("Message 1"),
				schema.NewAIChatMessage("Message 2"),
			}

			assert.Equal(t, expectedMessages, messages)
			mockClient.AssertExpectations(t)
		})

		t.Run("Messages returns an empty slice if there are no messages", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("LRange", mock.Anything, redisHistory.key(), int64(0), int64(-1)).
				Return(redis.Nil)

			// Call the Messages method
			messages, err := redisHistory.Messages()
			assert.NoError(t, err)

			// Assert that an empty slice is returned
			assert.Empty(t, messages)
			mockClient.AssertExpectations(t)
		})
	})

	t.Run("AddUserMessage", func(t *testing.T) {
		mockClient := &mockRedisClient{}
		redisHistory := NewRedis(mockClient, "session1")
		message := schema.NewHumanChatMessage("Hello, world!")

		t.Run("AddUserMessage adds the user message", func(t *testing.T) {
			redisMessage := schema.ChatMessageToMap(message)
			messageJSON, _ := json.Marshal(redisMessage)

			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), string(messageJSON)).
				Return(int64(1))

			err := redisHistory.AddUserMessage("Hello, world!")

			assert.NoError(t, err)
			mockClient.AssertExpectations(t)
		})

		t.Run("AddUserMessage returns an error if LPush fails", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), mock.Anything).
				Return(errors.New("LPush failed"))

			err := redisHistory.AddUserMessage("Hello, world!")

			assert.Error(t, err)
			mockClient.AssertExpectations(t)
		})
	})

	t.Run("AddAIMessage", func(t *testing.T) {
		mockClient := &mockRedisClient{}
		redisHistory := NewRedis(mockClient, "session1")
		message := schema.NewAIChatMessage("AI response")

		t.Run("AddAIMessage adds the AI message", func(t *testing.T) {
			redisMessage := schema.ChatMessageToMap(message)
			messageJSON, _ := json.Marshal(redisMessage)

			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), string(messageJSON)).
				Return(int64(1))

			err := redisHistory.AddAIMessage("AI response")

			assert.NoError(t, err)
			mockClient.AssertExpectations(t)
		})

		t.Run("AddAIMessage returns an error if LPush fails", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), mock.Anything).
				Return(errors.New("LPush failed"))

			err := redisHistory.AddAIMessage("AI response")

			assert.Error(t, err)
			mockClient.AssertExpectations(t)
		})
	})

	t.Run("AddMessage", func(t *testing.T) {
		mockClient := &mockRedisClient{}
		redisHistory := NewRedis(mockClient, "session1")
		message := schema.NewHumanChatMessage("Hello, world!")

		t.Run("AddMessage adds the chat message", func(t *testing.T) {
			redisMessage := schema.ChatMessageToMap(message)
			messageJSON, _ := json.Marshal(redisMessage)

			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), string(messageJSON)).
				Return(int64(1))

			err := redisHistory.AddMessage(message)

			assert.NoError(t, err)
			mockClient.AssertExpectations(t)
		})

		t.Run("AddMessage returns an error if LPush fails", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("LPush", mock.Anything, redisHistory.key(), mock.Anything).
				Return(errors.New("LPush failed"))

			err := redisHistory.AddMessage(message)

			assert.Error(t, err)
			mockClient.AssertExpectations(t)
		})
	})

	t.Run("Clear", func(t *testing.T) {
		mockClient := &mockRedisClient{}
		redisHistory := NewRedis(mockClient, "session1")

		t.Run("Clear deletes the chat message history", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("Del", mock.Anything, redisHistory.key()).
				Return(int64(1))

			err := redisHistory.Clear()
			assert.NoError(t, err)

			mockClient.AssertExpectations(t)
		})

		t.Run("Clear returns an error if Del fails", func(t *testing.T) {
			mockClient.Mock = mock.Mock{}
			mockClient.On("Del", mock.Anything, mock.Anything).
				Return(errors.New("del failed"))

			err := redisHistory.Clear()
			assert.Error(t, err)

			mockClient.AssertExpectations(t)
		})
	})
}
