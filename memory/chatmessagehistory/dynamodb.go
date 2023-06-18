package chatmessagehistory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hupe1980/golc/schema"
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
	History   schema.ChatMessages `dynamodbav:"history"`
}

// func (h *dynamoDBHistory) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
// 	fmt.Println("XXxxxXXXXXXXXXXXXXXXXXXXXX")
// 	fields, ok := value.
// 	if !ok {
// 		return nil
// 	}
// 	return nil
// }

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

func (mh *DynamoDB) Messages() (schema.ChatMessages, error) {
	sessionID, err := attributevalue.Marshal(mh.sessionID)
	if err != nil {
		return nil, err
	}

	result, err := mh.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"SessionId": sessionID,
		},
		TableName: aws.String(mh.tableName),
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return schema.ChatMessages{}, nil
	}

	out := make(map[string]interface{})
	if err := attributevalue.UnmarshalMap(result.Item, &out); err != nil {
		return nil, err
	}

	output := dynamoDBHistory{
		SessionID: out["SessionId"].(string),
		History:   nil,
	}

	return output.History, nil
}

func (mh *DynamoDB) AddUserMessage(text string) error {
	message := schema.NewHumanChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *DynamoDB) AddAIMessage(text string) error {
	message := schema.NewAIChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *DynamoDB) AddMessage(message schema.ChatMessage) error {
	messages, err := mh.Messages()
	if err != nil {
		return err
	}

	item, err := attributevalue.MarshalMap(dynamoDBHistory{
		SessionID: mh.sessionID,
		History:   append(messages, message),
	})
	if err != nil {
		return err
	}

	if _, err := mh.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(mh.tableName),
	}); err != nil {
		return err
	}

	return nil
}

func (mh *DynamoDB) Clear() error {
	sessionID, err := attributevalue.Marshal(mh.sessionID)
	if err != nil {
		return err
	}

	if _, err := mh.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"SessionId": sessionID,
		},
	}); err != nil {
		return err
	}

	return nil
}
