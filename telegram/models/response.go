package models

import (
	"encoding/json"
)

type Response struct {
	Message string `json:"message"`
}

type APIGatewayResponse struct {
	Body string `json:"body"`
	Headers map[string]string `json:"headers"`
	StatusCode int `json:"statusCode"`
}

func (r Response) GatewayResponse() (APIGatewayResponse) {
	body, _ := json.Marshal(r)
	return APIGatewayResponse{
		Body: string(body),
		Headers: make(map[string]string),
		StatusCode: 200,
	}
}
