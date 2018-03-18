package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/col/whosinbot/models"
)

func Handler() (APIGatewayResponse, error) {
	response := Response{
		Message: "WhosInBot",
	}

	return response.GatewayResponse(), nil
}

func main() {
	lambda.Start(Handler)
}
