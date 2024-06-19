package main_test

import (
	"guestbook/src/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type mockDynamoDBClient struct{}

func (m *mockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"Name":    {S: aws.String("John Doe")},
				"Message": {S: aws.String("Hello, world!")},
			},
		},
	}, nil
}

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{}, nil
}

func (m *mockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}

func TestGuestbookHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mockSvc := &mockDynamoDBClient{}
	h := handler.NewHandler(mockSvc, "GuestbookEntries", []string{"badword1", "badword2"})

	http.HandlerFunc(h.GuestbookHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Guestbook"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSignHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/sign", strings.NewReader("name=John&message=Hello"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	mockSvc := &mockDynamoDBClient{}
	h := handler.NewHandler(mockSvc, "GuestbookEntries", []string{"badword1", "badword2"})

	http.HandlerFunc(h.SignHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
}
