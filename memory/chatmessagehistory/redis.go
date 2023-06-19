package chatmessagehistory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hupe1980/golc/schema"
)

// Compile time check to ensure Redis satisfies the ChatMessageHistory interface.
var _ schema.ChatMessageHistory = (*Redis)(nil)

type RedisClient interface {
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

type RedisOptions struct {
	KeyPrefix string
	TTL       *time.Duration
}

type Redis struct {
	sessionID   string
	redisClient RedisClient
	opts        RedisOptions
}

func NewRedis(redisClient RedisClient, sessionID string) *Redis {
	opts := RedisOptions{
		KeyPrefix: "message_store:",
	}

	return &Redis{
		sessionID:   sessionID,
		redisClient: redisClient,
		opts:        opts,
	}
}

func (mh *Redis) Messages(ctx context.Context) (schema.ChatMessages, error) {
	messages := schema.ChatMessages{}

	items, err := mh.redisClient.LRange(ctx, mh.key(), 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return messages, nil
		}

		return nil, err
	}

	if len(items) == 0 {
		return messages, nil
	}

	for _, item := range items {
		message := map[string]string{}

		if err = json.Unmarshal([]byte(item), &message); err != nil {
			return nil, err
		}

		cm, err := schema.MapToChatMessage(message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, cm)
	}

	return messages, nil
}

func (mh *Redis) AddUserMessage(ctx context.Context, text string) error {
	message := schema.NewHumanChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *Redis) AddAIMessage(ctx context.Context, text string) error {
	message := schema.NewAIChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *Redis) AddMessage(ctx context.Context, message schema.ChatMessage) error {
	redisMessage := schema.ChatMessageToMap(message)

	messageJSON, err := json.Marshal(redisMessage)
	if err != nil {
		return err
	}

	if err := mh.redisClient.LPush(ctx, mh.key(), string(messageJSON)).Err(); err != nil {
		return err
	}

	if mh.opts.TTL != nil {
		if err := mh.redisClient.Expire(ctx, mh.key(), *mh.opts.TTL).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (mh *Redis) Clear(ctx context.Context) error {
	res := mh.redisClient.Del(ctx, mh.key())
	return res.Err()
}

func (mh *Redis) key() string {
	return mh.opts.KeyPrefix + mh.sessionID
}
