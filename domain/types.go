package domain

import "strings"
import "time"

type Command struct {
	ChatID string
	Name   string
	Params []string
	From   User
}

func (c Command) ParamsString() string {
	return strings.Join(c.Params, " ")
}

func (c Command) FirstParam() string {
	if len(c.Params) > 0 {
		return c.Params[0]
	}
	return ""
}

func (c Command) ParamsStringExceptFirst() string {
	if len(c.Params) > 1 {
		return strings.Join(c.Params[1:], " ")
	}
	return ""
}

type User struct {
	UserID string
	Name   string
}

func EmptyCommand() Command {
	return Command{
		ChatID: "",
		Name:   "",
		Params: nil,
		From:   User{},
	}
}

type Response struct {
	ChatID string
	Text   string
}

type DataStore interface {
	StartRollCall(rollCall RollCall) error
	EndRollCall(rollCall RollCall) error

	SetTitle(rollCall RollCall, title string) error
	SetQuiet(rollCall RollCall, quiet bool) error

	SetResponse(rollCallResponse RollCallResponse) error

	GetRollCall(chatID string) (*RollCall, error)
	LoadRollCallResponses(rollCall *RollCall) error
}

type RollCall struct {
	ChatID string
	Title  string
	Quiet  bool
	In     []RollCallResponse
	Out    []RollCallResponse
	Maybe  []RollCallResponse
}

func (r *RollCall) AddResponse(response RollCallResponse) {
	switch response.Status {
	case "in":
		r.In = append(r.In, response)
	case "out":
		r.Out = append(r.Out, response)
	case "maybe":
		r.Maybe = append(r.Maybe, response)
	}
}

type RollCallResponse struct {
	ChatID string
	UserID string
	Name   string
	Status string
	Reason string
	Date   time.Time
}

func NewRollCallResponse(command Command, name string, status string, reason string, isSetStatusFor bool) RollCallResponse {
	var userId = command.From.UserID
	if isSetStatusFor {
		userId = name
	}
	return RollCallResponse{
		ChatID: command.ChatID,
		UserID: userId,
		Name:   name,
		Status: status,
		Reason: reason,
		Date:   time.Now(),
	}
}

type Responses []RollCallResponse

func (r Responses) Len() int           { return len(r) }
func (r Responses) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Responses) Less(i, j int) bool { return r[i].Date.Before(r[j].Date) }
