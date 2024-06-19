package handler

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type MockDynamoDBClient struct {
	mu      sync.Mutex
	entries []GuestbookEntry
}

func (m *MockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	items := []map[string]*dynamodb.AttributeValue{}
	for _, entry := range m.entries {
		av, err := dynamodbattribute.MarshalMap(entry)
		if err != nil {
			return nil, err
		}
		items = append(items, av)
	}

	return &dynamodb.ScanOutput{
		Items: items,
	}, nil
}

func (m *MockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var entry GuestbookEntry
	if err := dynamodbattribute.UnmarshalMap(input.Item, &entry); err != nil {
		return nil, err
	}
	m.entries = append(m.entries, entry)

	return &dynamodb.PutItemOutput{}, nil
}

func (m *MockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{}, nil
}

func (m *MockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}
