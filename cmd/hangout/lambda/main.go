package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/chat/v1"
	"log"
	"whosinbot/hangout"
	"whosinbot/dynamodb"
	"whosinbot/whosinbot"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request Body: " + request.Body)

	command, err := hangout.ParseDeprecatedEvent([]byte(request.Body))
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
	if response == nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	message := chat.Message{
		Text: response.Text,
	}
	messageJOSN, _ := message.MarshalJSON()

	resBody := string(messageJOSN)
	log.Printf("response body: %s\n", resBody)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: resBody}, nil
}


