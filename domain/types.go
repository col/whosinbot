package domain

import "strings"

type Command struct {
	ChatID int64
	Name   string
	Params []string
	From   User
}

func (c Command) ParamsString() string {
	return strings.Join(c.Params, " ")
}

type User struct {
	UserID   int64
	Name string
}

func EmptyCommand() Command {
	return Command{
		ChatID: 0,
		Name:   "",
		Params: nil,
		From:   User{},
	}
}

type Response struct {
	ChatID int64
	Text   string
}

type DataStore interface {
	StartRollCall(rollCall RollCall) error
	EndRollCall(rollCall RollCall) error

	SetTitle(rollCall RollCall, title string) error
	SetQuiet(rollCall RollCall, quiet bool) error

	SetResponse(rollCallResponse RollCallResponse) error

	GetRollCall(chatID int64) (*RollCall, error)
	LoadRollCallResponses(rollCall *RollCall) error
}

type RollCall struct {
	ChatID int64
	Title  string
	Quiet  bool
	In     []RollCallResponse
	Out    []RollCallResponse
	Maybe  []RollCallResponse
}

func (r *RollCall)AddResponse(response RollCallResponse) {
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
	ChatID   int64
	UserID   int64
	Name     string
	Status   string
	Reason   string
}

func NewRollCallResponse(command Command, status string) RollCallResponse {
	return RollCallResponse{
		ChatID: command.ChatID,
		UserID: command.From.UserID,
		Name: command.From.Name,
		Status: status,
		Reason: command.ParamsString(),
	}
}