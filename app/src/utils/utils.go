package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func InitializeSession(region string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
}

func InitializeDynamoDB(sess *session.Session) *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func FetchForbiddenWords(url string) ([]string, error) {
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

func IsBadName(name string, forbiddenWords []string) (bool, string) {
	for _, badName := range forbiddenWords {
		if strings.EqualFold(name, badName) {
			return true, badName
		}
	}
	return false, ""
}
