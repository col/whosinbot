package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"whosinbot/domain"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

const EmptyString = " "

type DynamoDataStore struct {
	svc DynamoService
}

type DynamoService interface {
	UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
}

func NewDynamoDataStore() *DynamoDataStore {
	return &DynamoDataStore{
		svc: getService(),
	}
}

func (d DynamoDataStore) StartRollCall(rollCall domain.RollCall) error {
	// DEBUG
	log.Printf("StartRollCall: %+v\n", rollCall)

	title := rollCall.Title
	if len(title) == 0 {
		title = EmptyString
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(title),
			},
			":quiet": {
				BOOL: aws.Bool(rollCall.Quiet),
			},
		},
		TableName:        aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key:              getRollCallKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title, quiet = :quiet"),
	}

	_, err := d.svc.UpdateItem(input)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) GetRollCall(chatID int64) (*domain.RollCall, error) {
	// DEBUG
	log.Printf("GetRollCall: %+d\n", chatID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key:       getRollCallKey(chatID),
	}

	result, err := d.svc.GetItem(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}
	log.Printf("Result: %+v\n", result)

	title := *result.Item["title"].S
	if title == EmptyString {
		title = ""
	}

	rollCall := domain.RollCall{
		ChatID: chatID,
		Title:  title,
		Quiet:  *result.Item["quiet"].BOOL,
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
		TableName:              aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
	}

	result, err := d.svc.Query(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var responses = []domain.RollCallResponse{}
	for _, item := range result.Items {

		reason := *item["reason"].S
		if reason == EmptyString {
			reason = ""
		}

		date := time.Now()
		err = date.UnmarshalText([]byte(*item["date"].S))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		response := domain.RollCallResponse{
			ChatID: chatID,
			UserID: *item["user_id"].S,
			Name:   *item["name"].S,
			Status: *item["status"].S,
			Reason: reason,
			Date:   date,
		}
		responses = append(responses, response)
	}
	sort.Sort(domain.Responses(responses))

	return responses, nil
}

func (d DynamoDataStore) SetResponse(rollCallResponse domain.RollCallResponse) error {
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
	dateString, _ := rollCallResponse.Date.MarshalText()
	values[":date"] = &dynamodb.AttributeValue{
		S: aws.String(string(dateString)),
	}

	nameField := "name"
	statusField := "status"
	dateField := "date"
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#name":   &nameField,
			"#status": &statusField,
			"#date":   &dateField,
		},
		ExpressionAttributeValues: values,
		TableName:                 aws.String(os.Getenv("ROLLCALL_RESPONSE_TABLE")),
		Key:                       getResponseKey(rollCallResponse.ChatID, rollCallResponse.UserID),
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set #name = :name, #status = :status, reason = :reason, #date = :date"),
	}

	_, err := d.svc.UpdateItem(input)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) EndRollCall(rollCall domain.RollCall) error {
	// DEBUG
	log.Printf("EndRollCall: %+v\n", rollCall)

	responses, err := d.getRollCallResponses(rollCall.ChatID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	deleteResponseRequests := []*dynamodb.WriteRequest{}
	for _, response := range responses {
		writeRequest := &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: getResponseKey(rollCall.ChatID, response.UserID),
			},
		}
		deleteResponseRequests = append(deleteResponseRequests, writeRequest)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			os.Getenv("ROLLCALL_TABLE"): {
				&dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: getRollCallKey(rollCall.ChatID),
					},
				},
			},
			os.Getenv("ROLLCALL_RESPONSE_TABLE"): deleteResponseRequests,
		},
	}

	_, err = d.svc.BatchWriteItem(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) SetTitle(rollCall domain.RollCall, title string) error {
	// DEBUG
	log.Printf("SetTitle: %+v\n", rollCall)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(title),
			},
		},
		TableName:        aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key:              getRollCallKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title"),
	}

	_, err := d.svc.UpdateItem(input)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (d DynamoDataStore) SetQuiet(rollCall domain.RollCall, quiet bool) error {
	// DEBUG
	log.Printf("SetQuiet: %+v\n", rollCall)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":quiet": {
				BOOL: aws.Bool(quiet),
			},
		},
		TableName:        aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key:              getRollCallKey(rollCall.ChatID),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set quiet = :quiet"),
	}

	_, err := d.svc.UpdateItem(input)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func getRollCallKey(chatID int64) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"chat_id": {
			N: aws.String(strconv.Itoa(int(chatID))),
		},
	}
}

func getResponseKey(chatID int64, userID string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"chat_id": {
			N: aws.String(strconv.Itoa(int(chatID))),
		},
		"user_id": {
			S: aws.String(userID),
		},
	}
}

func getService() *dynamodb.DynamoDB {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	return dynamodb.New(sess)
}
