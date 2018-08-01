package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"whosinbot/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockDynamoService struct {
	UpdateItemCalled     bool
	UpdateItemCalledWith *dynamodb.UpdateItemInput
	UpdateItemResult     *dynamodb.UpdateItemOutput

	GetItemCalled     bool
	GetItemCalledWith *dynamodb.GetItemInput
	GetItemResult     *dynamodb.GetItemOutput

	QueryCalled     bool
	QueryCalledWith *dynamodb.QueryInput
	QueryResult     *dynamodb.QueryOutput

	BatchWriteItemCalled     bool
	BatchWriteItemCalledWith *dynamodb.BatchWriteItemInput
	BatchWriteItemResult     *dynamodb.BatchWriteItemOutput
}

func (d *MockDynamoService) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	d.UpdateItemCalled = true
	d.UpdateItemCalledWith = input
	return d.UpdateItemResult, nil
}

func (d *MockDynamoService) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	d.GetItemCalled = true
	d.GetItemCalledWith = input
	return d.GetItemResult, nil
}

func (d *MockDynamoService) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	d.QueryCalled = true
	d.QueryCalledWith = input
	return d.QueryResult, nil
}

func (d *MockDynamoService) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error) {
	d.BatchWriteItemCalled = true
	d.BatchWriteItemCalledWith = input
	return d.BatchWriteItemResult, nil
}

var mockDynamoService *MockDynamoService
var dataSource *DynamoDataStore

func setUp() {
	mockDynamoService = &MockDynamoService{
		UpdateItemCalled:         false,
		UpdateItemCalledWith:     nil,
		UpdateItemResult:         nil,
		GetItemCalled:            false,
		GetItemCalledWith:        nil,
		GetItemResult:            nil,
		QueryCalled:              false,
		QueryCalledWith:          nil,
		QueryResult:              nil,
		BatchWriteItemCalled:     false,
		BatchWriteItemCalledWith: nil,
		BatchWriteItemResult:     nil,
	}
	dataSource = &DynamoDataStore{svc: mockDynamoService}
}

func TestStartRollCall(t *testing.T) {
	setUp()
	rollCall := domain.RollCall{
		ChatID: 123,
		Title:  "Dinner",
		Quiet:  false,
	}
	dataSource.StartRollCall(rollCall)
	assert.True(t, mockDynamoService.UpdateItemCalled)
	assert.Equal(t, "Dinner", *mockDynamoService.UpdateItemCalledWith.ExpressionAttributeValues[":title"].S)
	assert.False(t, false, mockDynamoService.UpdateItemCalledWith.ExpressionAttributeValues[":quiet"].BOOL)
}

func TestStartRollCallSetsTitleToSpaceWhenEmpty(t *testing.T) {
	setUp()
	rollCall := domain.RollCall{
		ChatID: 123,
	}
	dataSource.StartRollCall(rollCall)
	assert.True(t, mockDynamoService.UpdateItemCalled)
	assert.Equal(t, " ", *mockDynamoService.UpdateItemCalledWith.ExpressionAttributeValues[":title"].S)
}
