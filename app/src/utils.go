package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func initializeSession(region string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
}

func initializeDynamoDB(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func fetchForbiddenWords(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	var words []string
	if err := json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, fmt.Errorf("failed to decode data: %v", err)
	}

	return words, nil
}

func isBadName(name string, forbiddenWords []string) bool {
	for _, word := range forbiddenWords {
		if strings.Contains(strings.ToLower(name), word) {
			return true
		}
	}
	return false
}

func isExplicitMessage(message string, forbiddenWords []string) bool {
	for _, word := range forbiddenWords {
		if strings.Contains(strings.ToLower(message), word) {
			return true
		}
	}
	return false
}
