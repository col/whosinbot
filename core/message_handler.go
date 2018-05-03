package core

import (
	"errors"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"os"
	"strconv"
)

func HandleMessage(message *tgbotapi.Message) (string, error) {
	switch command := message.Command(); command {
	case "start_roll_call":
		return handleStart(message)
	case "end_roll_call":
		return handleEnd(message)
	case "in":
		return handleIn(message)
	case "out":
		return handleOut(message)
	default:
		return "", errors.New("Not a bot command")
	}
}

func handleStart(message *tgbotapi.Message) (string, error) {
	// TODO: Create new roll call

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String("Testing"),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"chat_id": {
				N: aws.String(strconv.Itoa(int(message.Chat.ID))),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
	}

	return "Roll call started", nil
}

func handleEnd(message *tgbotapi.Message) (string, error) {
	// TODO: Delete roll call
	return "Roll call ended", nil
}

func handleIn(message *tgbotapi.Message) (string, error) {
	// TODO: mark sender as in
	return whosIn(message)
}

func handleOut(message *tgbotapi.Message) (string, error) {
	// TODO: mark sender as out
	return whosIn(message)
}

func whosIn(message *tgbotapi.Message) (string, error) {
	// TODO: output the roll call
	return "I don't know yet", nil
}
