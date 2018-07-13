package dynamodb

import (
	"github.com/col/whosinbot/domain"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strconv"
	"log"
)

type DynamoDataStore struct {

}

func (d DynamoDataStore) GetRollCall(chatID int64) (*domain.RollCall, error) {
	svc := getService()

	// DEBUG
	log.Printf("GetRollCall: %+d\n", chatID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: getKey(chatID),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	rollCall := domain.RollCall{
		ChatID: chatID,
		Title: *result.Item["title"].S,
		Quiet: *result.Item["quiet"].BOOL,
	}
	return &rollCall, nil
}

func (d DynamoDataStore) SetResponse(rollCallResponse domain.RollCallResponse) error {
	svc := getService()

	// DEBUG
	log.Printf("SetResponse: %+v\n", rollCallResponse)

	values := map[string]*dynamodb.AttributeValue{}
	values[":username"] = &dynamodb.AttributeValue{
		S: aws.String(rollCallResponse.Name),
	}
	values[":status"] = &dynamodb.AttributeValue{
		S: aws.String(rollCallResponse.Status),
	}
	if len(rollCallResponse.Reason) > 0 {
		values[":reason"] = &dynamodb.AttributeValue{
			S: aws.String(rollCallResponse.Reason),
		}
	}

	statusName := "status"
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{"#status": &statusName},
		ExpressionAttributeValues: values,
		TableName: aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
		Key: getResponseKey(rollCallResponse.ChatID, rollCallResponse.UserID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set username = :username, #status = :status"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) EndRollCall(rollCall domain.RollCall) error {

	// DEBUG
	log.Printf("EndRollCall: %+v\n", rollCall)

	// TODO: ...
	return nil
}

func (d DynamoDataStore) SetQuiet(rollCall domain.RollCall, quiet bool) error {
	svc := getService()

	// DEBUG
	log.Printf("SetQuiet: %+v\n", rollCall)

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
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) StartRollCall(rollCall domain.RollCall) (error) {
	svc := getService()

	// DEBUG
	log.Printf("StartRollCall: %+v\n", rollCall)

	var input *dynamodb.UpdateItemInput

	if len(rollCall.Title) > 0 {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":title": {
					S: aws.String(rollCall.Title),
				},
				":quiet": {
					BOOL: aws.Bool(rollCall.Quiet),
				},
			},
			TableName:        aws.String(os.Getenv("ROLLCALL_TABLE")),
			Key:              getKey(rollCall.ChatID),
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set title = :title, quiet = :quiet"),
		}
	} else {
		input = &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":quiet": {
					BOOL: aws.Bool(rollCall.Quiet),
				},
			},
			TableName:        aws.String(os.Getenv("ROLLCALL_TABLE")),
			Key:              getKey(rollCall.ChatID),
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set quiet = :quiet"),
		}
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		log.Println(err.Error())
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
