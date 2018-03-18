package main

import (
	"os"
	"fmt"
	"errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/telegram-bot-api.v4"
	"encoding/json"
)

func ValidateToken(requestToken string) (string, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token != requestToken {
		message := fmt.Sprintf("Bot token doesn't match! Expected: %v Received: %v", token, requestToken)
		return "", errors.New(message)
	}
	return token, nil
}

func ParseRequest(requestBody string) (*tgbotapi.Update, error) {
	update := &tgbotapi.Update{}
	err := json.Unmarshal([]byte(requestBody), update)
	return update, err
}

func HandleUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	bot.Send(msg)
}

func response(err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	} else {
		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token, err := ValidateToken(request.PathParameters["token"])
	if err != nil {
		return response(err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return response(err)
	}

	update, err := ParseRequest(request.Body)
	if err != nil {
		return response(err)
	}

	HandleUpdate(update, bot)

	return response(nil)
}

func main() {
	lambda.Start(Handler)
}
