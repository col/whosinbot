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
	Username string
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

	SetQuiet(rollCall RollCall, quiet bool) error

	SetResponse(rollCallResponse RollCallResponse) error

	GetRollCall(chatID int64) (*RollCall, error)
}

type RollCall struct {
	ChatID int64
	Title  string
	Quiet  bool
	In     []RollCallResponse
	Out    []RollCallResponse
	Maybe  []RollCallResponse
}

type RollCallResponse struct {
	ChatID   int64
	UserID   int64
	Name     string
	Response string
	Reason   string
}
