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
	// TODO: ...
	return nil, nil
}

func (d DynamoDataStore) SetResponse(rollCallResponse domain.RollCallResponse) (error) {
	// TODO: ...
	return nil
}

func (d DynamoDataStore) EndRollCall(rollCall domain.RollCall) error {
	// TODO: ...
	return nil
}

func (d DynamoDataStore) SetQuiet(rollCall domain.RollCall, quiet bool) error {
	// TODO: ...
	return nil
}

func (d DynamoDataStore) StartRollCall(rollCall domain.RollCall) (error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(rollCall.Title),
			},
		},
		TableName: aws.String(os.Getenv("ROLLCALL_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"chat_id": {
				N: aws.String(strconv.Itoa(int(rollCall.ChatID))),
			},
		},
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
