package main

import (
	"fmt"
	"whosinbot/dynamodb"
	whttp "whosinbot/http"
	"whosinbot/whosinbot"
	"net/http"
)

func main() {
	dataStore := dynamodb.NewDynamoDataStore()
	bot := &whosinbot.WhosInBot{DataStore: dataStore}

	http.Handle("/webhook", &whttp.WebhookHandler{WhosInBot: bot})
	port := 8080
	serverConfig := fmt.Sprintf(":%d", port)
	http.ListenAndServe(serverConfig, nil)
}
