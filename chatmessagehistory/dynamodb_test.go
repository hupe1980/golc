package chatmessagehistory

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hupe1980/golc/internal/util"
	"github.com/hupe1980/golc/schema"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDB_Messages(t *testing.T) {
	// Create a mock DynamoDB client
	mockClient := &mockDynamoDBClient{}

	// Create a test DynamoDB instance
	dynamoDB := NewDynamoDB(mockClient, "testTable", "testSessionID")

	t.Run("Messages returns history when it exists", func(t *testing.T) {
		msg1 := schema.NewHumanChatMessage("Message 1")
		msg2 := schema.NewAIChatMessage("Message 2")

		// Set up the mock client to return a valid response
		mockClient.GetItemFunc = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			item, err := attributevalue.MarshalMap(dynamoDBHistory{
				SessionID: "testSessionID",
				History: util.Map(schema.ChatMessages{msg1, msg2}, func(m schema.ChatMessage, _ int) map[string]string {
					return schema.ChatMessageToMap(m)
				}),
			})

			return &dynamodb.GetItemOutput{Item: item}, err
		}

		// Call the Messages method
		messages, err := dynamoDB.Messages(context.TODO())

		// Assert that the expected history is returned
		expectedHistory := schema.ChatMessages{msg1, msg2}

		assert.NoError(t, err)
		assert.Equal(t, expectedHistory, messages)
	})

	t.Run("Messages returns empty history when it does not exist", func(t *testing.T) {
		// Set up the mock client to return an empty response
		mockClient.GetItemFunc = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{}, nil
		}

		// Call the Messages method
		messages, err := dynamoDB.Messages(context.TODO())

		// Assert that an empty history is returned
		assert.NoError(t, err)
		assert.Empty(t, messages)
	})
}

// Mock DynamoDB client implementation
type mockDynamoDBClient struct {
	dynamodb.Client
	GetItemFunc func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func (m *mockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.GetItemFunc != nil {
		return m.GetItemFunc(ctx, params, optFns...)
	}

	return nil, nil
}
