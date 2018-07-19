package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/col/whosinbot/dynamodb"
	"github.com/col/whosinbot/telegram"
	"github.com/col/whosinbot/whosinbot"
	"log"
	"strings"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// DEBUG
	log.Printf("Request Body: " + request.Body)

	if strings.Contains(request.Body, "edited_message") {
		// Avoid crash caused by the telegram lib not handling edited messages correctly yet.
		// Just ignore this message.
		log.Printf("Ignoring edited message. Not supported")
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	// Parse Command
	command, err := telegram.ParseUpdate([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	// Process Command
	dataStore := dynamodb.NewDynamoDataStore()
	bot := whosinbot.WhosInBot{DataStore: dataStore}
	response, err := bot.HandleCommand(command)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	// Send Response
	if response != nil {
		api := telegram.NewTelegram(request.PathParameters["token"])
		err = api.SendResponse(response)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400}, err
		}
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}
