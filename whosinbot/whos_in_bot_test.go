package whosinbot

import (
	"whosinbot/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockDataStore struct {
	startRollCallCalled bool
	startRollCallWith   *domain.RollCall

	endRollCallCalled bool
	endRollCallWith   *domain.RollCall

	setTitleCalled       bool
	setTitleWithRollCall *domain.RollCall
	setTitleWithTitle    string

	setQuietCalled       bool
	setQuietWithRollCall *domain.RollCall
	setQuietWithBool     bool

	setResponseCalled bool
	setResponseWith   *domain.RollCallResponse

	getRollCallCalled bool
	getRollCallWith   string

	loadRollCallResponsesCalled bool
	loadRollCallResponsesWith   *domain.RollCall

	rollCall          *domain.RollCall
	rollCallResponses []domain.RollCallResponse
}

func (d *MockDataStore) StartRollCall(rollCall domain.RollCall) error {
	d.startRollCallCalled = true
	d.startRollCallWith = &rollCall
	return nil
}

func (d *MockDataStore) EndRollCall(rollCall domain.RollCall) error {
	d.endRollCallCalled = true
	d.endRollCallWith = &rollCall
	return nil
}

func (d *MockDataStore) SetTitle(rollCall domain.RollCall, title string) error {
	d.setTitleCalled = true
	d.setTitleWithRollCall = &rollCall
	d.setTitleWithTitle = title
	return nil
}

func (d *MockDataStore) SetQuiet(rollCall domain.RollCall, quiet bool) error {
	d.setQuietCalled = true
	d.setQuietWithRollCall = &rollCall
	d.setQuietWithBool = quiet
	return nil
}

func (d *MockDataStore) SetResponse(rollCallResponse domain.RollCallResponse) error {
	d.setResponseCalled = true
	d.setResponseWith = &rollCallResponse
	d.rollCallResponses = append(d.rollCallResponses, rollCallResponse)
	return nil
}

func (d *MockDataStore) GetRollCall(chatID string) (*domain.RollCall, error) {
	d.getRollCallCalled = true
	d.getRollCallWith = chatID
	return d.rollCall, nil
}

func (d *MockDataStore) LoadRollCallResponses(rollCall *domain.RollCall) error {
	d.loadRollCallResponsesCalled = true
	d.loadRollCallResponsesWith = rollCall
	for _, response := range d.rollCallResponses {
		rollCall.AddResponse(response)
	}
	return nil
}

var mockDataStore *MockDataStore
var bot *WhosInBot

func setUp() {
	mockDataStore = &MockDataStore{
		startRollCallCalled: false,
		startRollCallWith:   nil,
		endRollCallCalled:   false,
		endRollCallWith:     nil,
		setResponseCalled:   false,
		setResponseWith:     nil,
		rollCall:            nil,
		rollCallResponses:   nil,
	}
	bot = &WhosInBot{DataStore: mockDataStore}
}

func TestStartRollCall(t *testing.T) {
	setUp()
	command := domain.Command{ChatID: "123", Name: "start_roll_call", Params: []string{"sample title"}}
	response, err := bot.HandleCommand(command)
	// Validate end previous roll call (if exists)
	assert.True(t, mockDataStore.endRollCallCalled)
	assert.Equal(t, "123", mockDataStore.endRollCallWith.ChatID)
	// Validate start new roll call
	assert.True(t, mockDataStore.startRollCallCalled)
	assert.NotNil(t, mockDataStore.startRollCallWith)
	assert.Equal(t, "123", mockDataStore.startRollCallWith.ChatID)
	assert.Equal(t, "sample title", mockDataStore.startRollCallWith.Title)
	// Validate response
	assertBotResponse(t, response, err, "123", "Roll call started", nil)
}

func TestEndRollCallWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	command := domain.Command{ChatID: "123", Name: "end_roll_call", Params: []string{}}
	response, err := bot.HandleCommand(command)
	// Validate data store
	assert.True(t, mockDataStore.endRollCallCalled)
	assert.NotNil(t, mockDataStore.endRollCallWith)
	assert.Equal(t, "123", mockDataStore.endRollCallWith.ChatID)
	// Validate response
	assertBotResponse(t, response, err, "123", "Roll call ended", nil)
}

func TestEndRollCallWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = nil
	command := domain.Command{ChatID: "123", Name: "end_roll_call", Params: []string{}}
	response, err := bot.HandleCommand(command)
	// Validate data store
	assert.False(t, mockDataStore.endRollCallCalled)
	// Validate response
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestSetTitleWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	response, err := bot.HandleCommand(command("set_title", []string{"foo", "bar"}))
	// Validate data store
	assert.True(t, mockDataStore.setTitleCalled, "should call setTitle")
	assert.NotNil(t, mockDataStore.setTitleWithRollCall)
	assert.Equal(t, "123", mockDataStore.setTitleWithRollCall.ChatID)
	assert.Equal(t, "foo bar", mockDataStore.setTitleWithTitle)
	// Validate response
	assertBotResponse(t, response, err, "123", "Roll call title set", nil)
}

func TestSetTitleWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("set_title", nil))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestInWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(responseCommand("in", []string{}))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestInWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}

	response, err := bot.HandleCommand(responseCommand("in", []string{}))

	assertResponsePersisted(t, "123", "456", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "1. JohnSmith", nil)
}

func TestInWithReasonWhenWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}

	response, err := bot.HandleCommand(responseCommand("in", []string{"sample", "reason"}))

	assertResponsePersisted(t, "123", "456", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "1. JohnSmith (sample reason)", nil)
}

func TestInWithReasonWithParenthesesWhenWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}

	response, err := bot.HandleCommand(responseCommand("in", []string{"(sample", "reason)"}))

	assertResponsePersisted(t, "123", "456", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "1. JohnSmith (sample reason)", nil)
}

func TestInWhenRollCallIsInQuietMode(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{
		ChatID: "123",
		Title:  "",
		Quiet:  true,
		In: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 1", Status: "in", Reason: ""},
		},
		Out: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 2", Status: "out", Reason: ""},
		},
		Maybe: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 3", Status: "maybe", Reason: ""},
		},
	}

	response, err := bot.HandleCommand(responseCommand("in", []string{}))
	assertResponsePersisted(t, "123", "456", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "JohnSmith is in!\nTotal: 2 in, 1 out, 1 might come\n", nil)
}

func TestOutWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(responseCommand("out", []string{}))
	assert.False(t, mockDataStore.setResponseCalled)
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestOutWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	response, err := bot.HandleCommand(responseCommand("out", []string{}))
	assertResponsePersisted(t, "123", "456", "out", "JohnSmith")
	assertBotResponse(t, response, err, "123", "Out\n - JohnSmith", nil)
}

func TestOutWithReasonWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	response, err := bot.HandleCommand(responseCommand("out", []string{"sample", "reason"}))
	assertResponsePersisted(t, "123", "456", "out", "JohnSmith")
	assertBotResponse(t, response, err, "123", "Out\n - JohnSmith (sample reason)", nil)
}

func TestMaybeWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(responseCommand("maybe", []string{}))
	assert.False(t, mockDataStore.setResponseCalled)
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestMaybeWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	response, err := bot.HandleCommand(responseCommand("maybe", []string{}))
	assertResponsePersisted(t, "123", "456", "maybe", "JohnSmith")
	assertBotResponse(t, response, err, "123", "Maybe\n - JohnSmith", nil)
}

func TestMaybeWithReasonWhenRollCallExists(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: ""}
	response, err := bot.HandleCommand(responseCommand("maybe", []string{"sample", "reason"}))
	assertResponsePersisted(t, "123", "456", "maybe", "JohnSmith")
	assertBotResponse(t, response, err, "123", "Maybe\n - JohnSmith (sample reason)", nil)
}

func TestWhosIn(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{
		ChatID: "123",
		Title:  "Test Title",
		In: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 1", Status: "in", Reason: ""},
		},
		Out: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 2", Status: "out", Reason: ""},
		},
		Maybe: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 3", Status: "maybe", Reason: ""},
		},
	}
	response, err := bot.HandleCommand(responseCommand("whos_in", []string{}))
	assertBotResponse(t, response, err, "123", "Test Title\n1. User 1\n\nOut\n - User 2\n\nMaybe\n - User 3", nil)
}

func TestWhosInWhenRollCallDoesNotExist(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(responseCommand("whos_in", []string{}))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestWhosInWhenThereAreNoResponses(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123", Title: "Test Title"}
	response, err := bot.HandleCommand(responseCommand("whos_in", []string{}))
	assertBotResponse(t, response, err, "123", "No responses yet. 😢", nil)
}

func TestShh(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{
		ChatID: "123",
		Title:  "Test Title",
		Quiet:  false,
	}

	response, err := bot.HandleCommand(command("shh", nil))
	// Validate data store
	assert.True(t, mockDataStore.setQuietCalled)
	assert.NotNil(t, mockDataStore.setQuietWithRollCall)
	assert.True(t, mockDataStore.setQuietWithBool)
	// Validate Response
	assertBotResponse(t, response, err, "123", "Ok fine, I'll be quiet. 🤐", nil)
}

func TestShhWhenRollCallDoesNotExists(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("shh", nil))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestLouder(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{
		ChatID: "123",
		In: []domain.RollCallResponse{
			{ChatID: "123", UserID: "1", Name: "User 1", Status: "in", Reason: ""},
		},
	}

	response, err := bot.HandleCommand(command("louder", nil))
	// Validate data store
	assert.True(t, mockDataStore.setQuietCalled)
	assert.NotNil(t, mockDataStore.setQuietWithRollCall)
	assert.False(t, mockDataStore.setQuietWithBool)
	// Validate Response
	assertBotResponse(t, response, err, "123", "Sure. 😃\n1. User 1", nil)
}

func TestLouderWhenRollCallDoesNotExist(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("louder", nil))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestSetInForWhenRollCallDoesNotExist(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("set_in_for", []string{"JohnSmith"}))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestSetInForWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_in_for", []string{"JohnSmith"}))

	assertResponsePersisted(t, "123", "JohnSmith", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "1. JohnSmith", nil)
}

func TestSetInForWithReasonWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_in_for", []string{"JohnSmith", "sample", "reason"}))

	assertResponsePersisted(t, "123", "JohnSmith", "in", "JohnSmith")
	assertBotResponse(t, response, err, "123", "1. JohnSmith (sample reason)", nil)
}

func TestSetInForWithoutANameWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_in_for", []string{}))
	assertBotResponse(t, response, err, "123", "Please provide the persons first name.", nil)
}

func TestSetOutForWhenRollCallDoesNotExist(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("set_out_for", []string{"JohnSmith"}))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestSetOutForWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_out_for", []string{"JohnSmith"}))
	assertBotResponse(t, response, err, "123", "Out\n - JohnSmith", nil)
}

func TestSetOutForWithReasonWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_out_for", []string{"JohnSmith", "sample", "reason"}))
	assertBotResponse(t, response, err, "123", "Out\n - JohnSmith (sample reason)", nil)
}

func TestSetOutForWithoutANameWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_out_for", []string{}))
	assertBotResponse(t, response, err, "123", "Please provide the persons first name.", nil)
}

func TestSetMaybeForWhenRollCallDoesNotExist(t *testing.T) {
	setUp()
	response, err := bot.HandleCommand(command("set_maybe_for", []string{"JohnSmith"}))
	assertBotResponse(t, response, err, "123", "No roll call in progress", nil)
}

func TestSetMaybeForWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_maybe_for", []string{"JohnSmith"}))
	assertBotResponse(t, response, err, "123", "Maybe\n - JohnSmith", nil)
}

func TestSetMaybeForWithReasonWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_maybe_for", []string{"JohnSmith", "sample", "reason"}))
	assertBotResponse(t, response, err, "123", "Maybe\n - JohnSmith (sample reason)", nil)
}

func TestSetMaybeForWithoutANameWhenRollCallExist(t *testing.T) {
	setUp()
	mockDataStore.rollCall = &domain.RollCall{ChatID: "123"}
	response, err := bot.HandleCommand(command("set_maybe_for", []string{}))
	assertBotResponse(t, response, err, "123", "Please provide the persons first name.", nil)
}

// Test Helpers
func command(name string, params []string) domain.Command {
	return domain.Command{
		ChatID: "123",
		Name:   name,
		Params: params,
		From:   domain.User{UserID: "456", Name: "JohnSmith"},
	}
}

func responseCommand(status string, params []string) domain.Command {
	return domain.Command{
		ChatID: "123",
		Name:   status,
		Params: params,
		From:   domain.User{UserID: "456", Name: "JohnSmith"},
	}
}

//func assertResponseNotPersisted(t *testing.T, chatID int, userID string) {
//	assert.NotNil(t, mockDataStore.setResponseWith)
//	if mockDataStore.setResponseWith != nil {
//		assert.Equal(t, int64(ch, mockDataStore.setResponseWith.ChatID)"
//		assert.Equal(t, userID, mockDataStore.setResponseWith.UserID)
//	}
//}


func assertResponsePersisted(t *testing.T, chatID string, userID string, status string, name string) {
	assert.True(t, mockDataStore.setResponseCalled, "should call setResponse")
	assert.NotNil(t, mockDataStore.setResponseWith)
	if mockDataStore.setResponseWith != nil {
		assert.Equal(t, chatID, mockDataStore.setResponseWith.ChatID)
		assert.Equal(t, userID, mockDataStore.setResponseWith.UserID)
		assert.Equal(t, status, mockDataStore.setResponseWith.Status)
		assert.Equal(t, name, mockDataStore.setResponseWith.Name)
	}
}

func assertBotResponse(t *testing.T, response *domain.Response, err error, chatID string, text string, error error) {
	assert.Equal(t, chatID, response.ChatID)
	assert.Equal(t, text, response.Text)
	assert.Equal(t, error, err)
}
