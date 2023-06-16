package chatmessagehistory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hupe1980/golc"
)

// Compile time check to ensure DynamoDB satisfies the ChatMessageHistory interface.
var _ golc.ChatMessageHistory = (*DynamoDB)(nil)

type dynamoDBHistory struct {
	SessionID string             `dynamodbav:"SessionId"`
	History   []golc.ChatMessage `dynamodbav:"History"`
}

type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
	sessionID string
}

func NewDynamoDB(client *dynamodb.Client, tableName, sessionID string) *DynamoDB {
	return &DynamoDB{
		client:    client,
		tableName: tableName,
		sessionID: sessionID,
	}
}

func (mh *DynamoDB) Messages() ([]golc.ChatMessage, error) {
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
		return []golc.ChatMessage{}, nil
	}

	output := dynamoDBHistory{}
	if err := attributevalue.UnmarshalMap(result.Item, &output); err != nil {
		return nil, err
	}

	return output.History, nil
}

func (mh *DynamoDB) AddUserMessage(text string) error {
	message := golc.NewHumanChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *DynamoDB) AddAIMessage(text string) error {
	message := golc.NewAIChatMessage(text)
	return mh.AddMessage(message)
}

func (mh *DynamoDB) AddMessage(message golc.ChatMessage) error {
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
