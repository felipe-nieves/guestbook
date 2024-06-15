// src/handler.go
package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBClient interface {
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

type GuestbookEntry struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Handler struct {
	svc            DynamoDBClient
	tableName      string
	forbiddenWords []string
	tpl            *template.Template
	mu             sync.Mutex
}

func NewHandler(svc DynamoDBClient, tableName string, forbiddenWords []string) *Handler {
	return &Handler{
		svc:            svc,
		tableName:      tableName,
		forbiddenWords: forbiddenWords,
		tpl: template.Must(template.New("guestbook").Parse(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>Guestbook</title>
            </head>
            <body>
                <h1>Guestbook</h1>
                <form action="/sign" method="POST">
                    <label for="name">Name:</label>
                    <input type="text" id="name" name="name" required>
                    <label for="message">Message:</label>
                    <textarea id="message" name="message" required></textarea>
                    <button type="submit">Sign</button>
                </form>
                <h2>Entries:</h2>
                <ul>
                    {{range .}}
                        <li><strong>{{.Name}}:</strong> {{.Message}}</li>
                    {{end}}
                </ul>
            </body>
            </html>
        `)),
	}
}

func (h *Handler) GuestbookHandler(w http.ResponseWriter, r *http.Request) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(h.tableName),
	}

	result, err := h.svc.Scan(scanInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var entries []GuestbookEntry
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.tpl.Execute(w, entries)
}

func (h *Handler) SignHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	message := r.FormValue("message")

	if isBadName(name, h.forbiddenWords) || isExplicitMessage(message, h.forbiddenWords) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	entry := GuestbookEntry{Name: name, Message: message}
	av, err := dynamodbattribute.MarshalMap(entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(h.tableName),
	}

	_, err = h.svc.PutItem(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
