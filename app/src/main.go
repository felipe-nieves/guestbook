package main

import (
	"net/http"
	"os"

	"guestbook/src/handler"
	"guestbook/src/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	region := os.Getenv("AWS_REGION")
	tableName := os.Getenv("TABLE_NAME")

	var svc handler.DynamoDBClient

	if region == "local" {
		svc = &handler.MockDynamoDBClient{}
		println("Using mock DynamoDB client")
	} else {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
		svc = dynamodb.New(sess)
		println("Using real DynamoDB client")
	}

	forbiddenWords, err := utils.FetchForbiddenWords("https://raw.githubusercontent.com/zacanger/profane-words/master/words.json")
	if err != nil {
		panic(err)
	}

	h := handler.NewHandler(svc, tableName, forbiddenWords)

	http.HandleFunc("/", h.GuestbookHandler)
	http.HandleFunc("/sign", h.SignHandler)
	http.ListenAndServe(":8080", nil)
}
