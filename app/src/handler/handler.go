package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"guestbook/src/utils"

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
	tpl, err := template.ParseFiles("src/html/index.html")
	if err != nil {
		panic(err)
	}
	return &Handler{
		svc:            svc,
		tableName:      tableName,
		forbiddenWords: forbiddenWords,
		tpl:            tpl,
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

	fmt.Printf("Received name: %s, message: %s\n", name, message)

	if isBad, badWord := utils.IsBadName(name, h.forbiddenWords); isBad {
		fmt.Printf("Name is flagged as bad: %s\n", badWord)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if isExplicit, explicitWord := h.isExplicitMessage(message); isExplicit {
		fmt.Printf("Message is flagged as explicit: %s\n", explicitWord)
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

func (h *Handler) isExplicitMessage(message string) (bool, string) {
	for _, word := range h.forbiddenWords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if re.MatchString(strings.ToLower(message)) {
			return true, word
		}
	}
	return false, ""
}
