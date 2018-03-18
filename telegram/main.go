package main

import (
	"context"
	"whosinbot/telegram/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := models.Response{
		Message: "WhosInBot - " + request.PathParameters["token"],
	}

	return events.APIGatewayProxyResponse{
		Body: response.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
