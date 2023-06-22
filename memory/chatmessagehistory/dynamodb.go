package chatmessagehistory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hupe1980/golc/schema"
	"github.com/hupe1980/golc/util"
)

// Compile time check to ensure DynamoDB satisfies the ChatMessageHistory interface.
var _ schema.ChatMessageHistory = (*DynamoDB)(nil)

type DynamoDBClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type dynamoDBHistory struct {
	SessionID string              `dynamodbav:"sessionId"`
	History   []map[string]string `dynamodbav:"history"`
}

type DynamoDB struct {
	client    DynamoDBClient
	tableName string
	sessionID string
}

func NewDynamoDB(client DynamoDBClient, tableName, sessionID string) *DynamoDB {
	return &DynamoDB{
		client:    client,
		tableName: tableName,
		sessionID: sessionID,
	}
}

func (mh *DynamoDB) Messages(ctx context.Context) (schema.ChatMessages, error) {
	sessionID, err := attributevalue.Marshal(mh.sessionID)
	if err != nil {
		return nil, err
	}

	result, err := mh.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"sessionId": sessionID,
		},
		TableName: aws.String(mh.tableName),
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return schema.ChatMessages{}, nil
	}

	output := dynamoDBHistory{}
	if err := attributevalue.UnmarshalMap(result.Item, &output); err != nil {
		return nil, err
	}

	history := schema.ChatMessages{}

	for _, v := range output.History {
		cm, err := schema.MapToChatMessage(v)
		if err != nil {
			return nil, err
		}

		history = append(history, cm)
	}

	return history, nil
}

func (mh *DynamoDB) AddUserMessage(ctx context.Context, text string) error {
	message := schema.NewHumanChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *DynamoDB) AddAIMessage(ctx context.Context, text string) error {
	message := schema.NewAIChatMessage(text)
	return mh.AddMessage(ctx, message)
}

func (mh *DynamoDB) AddMessage(ctx context.Context, message schema.ChatMessage) error {
	messages, err := mh.Messages(ctx)
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(dynamoDBHistory{
		SessionID: mh.sessionID,
		History: util.Map(append(messages, message), func(m schema.ChatMessage, _ int) map[string]string {
			return schema.ChatMessageToMap(m)
		}),
	})
	if err != nil {
		return err
	}

	if _, err := mh.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(mh.tableName),
	}); err != nil {
		return err
	}

	return nil
}

func (mh *DynamoDB) Clear(ctx context.Context) error {
	sessionID, err := attributevalue.Marshal(mh.sessionID)
	if err != nil {
		return err
	}

	if _, err := mh.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"sessionId": sessionID,
		},
	}); err != nil {
		return err
	}

	return nil
}
