package main

import (
	"whosinbot/telegram/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (models.APIGatewayResponse, error) {
	response := models.Response{
		Message: "WhosInBot",
	}

	return response.GatewayResponse(), nil
}

func main() {
	lambda.Start(Handler)
}
