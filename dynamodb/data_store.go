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

const EmptyString = " "

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

func (d DynamoDataStore) LoadRollCallResponses(rollCall *domain.RollCall) error {
	responses, err := d.getRollCallResponses(rollCall.ChatID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, response := range responses {
		rollCall.AddResponse(response)
	}

	return nil
}

func (d DynamoDataStore) getRollCallResponses(chatID int64) ([]domain.RollCallResponse, error) {
	svc := getService()

	// DEBUG
	log.Printf("getRollCallResponses ChatID: %+d\n", chatID)

	keyConditionExpression := "chat_id = :chat_id"
	selectString := dynamodb.SelectAllAttributes

	input := &dynamodb.QueryInput{
		Select: &selectString,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":chat_id": {
				N: aws.String(strconv.Itoa(int(chatID))),
			},
		},
		KeyConditionExpression: &keyConditionExpression,
		TableName: aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
	}

	result, err := svc.Query(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var responses = []domain.RollCallResponse{}
	for _, item := range result.Items {
		userID, err := strconv.Atoi(*item["user_id"].N)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		reason := *item["reason"].S
		if reason == EmptyString {
			reason = ""
		}

		response := domain.RollCallResponse{
			ChatID: chatID,
			UserID: int64(userID),
			Name: *item["name"].S,
			Status: *item["status"].S,
			Reason: reason,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (d DynamoDataStore) SetResponse(rollCallResponse domain.RollCallResponse) error {
	svc := getService()

	// DEBUG
	log.Printf("SetResponse: %+v\n", rollCallResponse)

	values := map[string]*dynamodb.AttributeValue{}
	values[":name"] = &dynamodb.AttributeValue{
		S: aws.String(rollCallResponse.Name),
	}
	values[":status"] = &dynamodb.AttributeValue{
		S: aws.String(rollCallResponse.Status),
	}
	reason := rollCallResponse.Reason
	if len(reason) == 0 {
		reason = EmptyString
	}
	values[":reason"] = &dynamodb.AttributeValue{
		S: aws.String(reason),
	}

	nameField := "name"
	statusField := "status"
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#name": &nameField,
			"#status": &statusField,
		},
		ExpressionAttributeValues: values,
		TableName: aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
		Key: getResponseKey(rollCallResponse.ChatID, rollCallResponse.UserID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set #name = :name, #status = :status, reason = :reason"),
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

func (d DynamoDataStore) SetTitle(rollCall domain.RollCall, title string) error {
	svc := getService()

	// DEBUG
	log.Printf("SetTitle: %+v\n", rollCall)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(title),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: getKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

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
