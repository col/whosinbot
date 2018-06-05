package main

import (
	"github.com/col/whosinbot/dynamodb"
	"net/http"
	"fmt"
	"github.com/col/whosinbot/whosinbot"
	whttp "github.com/col/whosinbot/http"
)

func main() {
	dataStore := &dynamodb.DynamoDataStore{}
	bot := &whosinbot.WhosInBot{ DataStore: dataStore }

	http.Handle("/webhook", &whttp.WebhookHandler{WhosInBot: bot})
	port := 8080
	serverConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(serverConfig, nil)
}

