package dynamodb

import (
	"github.com/col/whosinbot/domain"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strconv"
	"fmt"
)

type DynamoDataStore struct {

}

func (d DynamoDataStore) GetRollCall(chatID int64) (*domain.RollCall, error) {
	svc := getService()

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: getKey(chatID),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	rollCall := domain.RollCall{
		ChatID: chatID,
		Title: result.Item["Title"].String(),
		Quiet: *result.Item["Quiet"].BOOL,
	}
	return &rollCall, nil
}

func (d DynamoDataStore) SetResponse(rollCallResponse domain.RollCallResponse) error {
	svc := getService()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(rollCallResponse.Name),
			},
			":response": {
				S: aws.String(rollCallResponse.Response),
			},
			":reason": {
				S: aws.String(rollCallResponse.Reason),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
		Key: getResponseKey(rollCallResponse.ChatID, rollCallResponse.UserID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set name = :name, response = :response, reason = :reason"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) EndRollCall(rollCall domain.RollCall) error {
	// TODO: ...
	return nil
}

func (d DynamoDataStore) SetQuiet(rollCall domain.RollCall, quiet bool) error {
	svc := getService()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":quiet": {
				BOOL: aws.Bool(quiet),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: getKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set quiet = :quiet"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) StartRollCall(rollCall domain.RollCall) (error) {
	svc := getService()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(rollCall.Title),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: getKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func getKey(chatID int64) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"chat_id": {
			N: aws.String(strconv.Itoa(int(chatID))),
		},
	}
}

func getResponseKey(chatID int64, userID int64) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"chat_id": {
			N: aws.String(strconv.Itoa(int(chatID))),
		},
		"user_id": {
			N: aws.String(strconv.Itoa(int(userID))),
		},
	}
}

func getService() *dynamodb.DynamoDB {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	return dynamodb.New(sess)
}
