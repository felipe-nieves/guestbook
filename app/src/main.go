package main

import (
	"net/http"
	"os"
)

func main() {
	region := os.Getenv("AWS_REGION")
	tableName := os.Getenv("TABLE_NAME")

	sess := initializeSession(region)
	svc := initializeDynamoDB(sess)

	forbiddenWords, err := fetchForbiddenWords("https://raw.githubusercontent.com/zacanger/profane-words/master/words.json")
	if err != nil {
		panic(err)
	}

	// Initialize handlers
	handler := NewHandler(svc, tableName, forbiddenWords)

	http.HandleFunc("/", handler.GuestbookHandler)
	http.HandleFunc("/sign", handler.SignHandler)
	http.ListenAndServe(":8080", nil)
}
