package whosinbot

import (
	"fmt"
	"github.com/col/whosinbot/domain"
	"log"
	"strings"
)

type WhosInBot struct {
	DataStore domain.DataStore
}

func (b *WhosInBot) HandleCommand(command domain.Command) (*domain.Response, error) {
	// DEBUG
	log.Printf("Command: %+v\n", command)

	switch command.Name {
	case "start_roll_call":
		return b.handleStart(command)
	case "end_roll_call":
		return b.handleEnd(command)
	case "set_title":
		return b.handleSetTitle(command)
	case "in":
		return b.handleResponse(command, "in")
	case "out":
		return b.handleResponse(command, "out")
	case "maybe":
		return b.handleResponse(command, "maybe")
	case "set_in_for":
		return b.handleResponseFor(command, "in")
	case "set_out_for":
		return b.handleResponseFor(command, "out")
	case "set_maybe_for":
		return b.handleResponseFor(command, "maybe")
	case "whos_in":
		return b.handleWhosIn(command)
	case "shh":
		return b.handleSetQuiet(command, true)
	case "louder":
		return b.handleSetQuiet(command, false)
	default:
		log.Printf("Not a bot command: %+v\n", command)
		return nil, nil
	}
}

func (b *WhosInBot) handleStart(command domain.Command) (*domain.Response, error) {
	rollCall := domain.RollCall{
		ChatID: command.ChatID,
		Title:  command.ParamsString(),
	}

	_ = b.DataStore.EndRollCall(rollCall)
	err := b.DataStore.StartRollCall(rollCall)
	if err != nil {
		return nil, err
	}

	return &domain.Response{ChatID: command.ChatID, Text: "Roll call started"}, nil
}

func (b *WhosInBot) handleEnd(command domain.Command) (*domain.Response, error) {
	rollCall, err := b.DataStore.GetRollCall(command.ChatID)
	if err != nil {
		return nil, err
	}
	if rollCall == nil {
		return &domain.Response{Text: "No roll call in progress", ChatID: command.ChatID}, nil
	}
	err = b.DataStore.EndRollCall(*rollCall)
	if err != nil {
		return nil, err
	}
	return &domain.Response{ChatID: command.ChatID, Text: "Roll call ended"}, nil
}

func (b *WhosInBot) handleSetTitle(command domain.Command) (*domain.Response, error) {
	rollCall, err := b.DataStore.GetRollCall(command.ChatID)
	if err != nil {
		return nil, err
	}
	if rollCall == nil {
		return &domain.Response{Text: "No roll call in progress", ChatID: command.ChatID}, nil
	}
	err = b.DataStore.SetTitle(*rollCall, command.ParamsString())
	if err != nil {
		return nil, err
	}
	return &domain.Response{ChatID: command.ChatID, Text: "Roll call title set"}, nil
}

func (b *WhosInBot) handleSetQuiet(command domain.Command, quiet bool) (*domain.Response, error) {
	rollCall, err := b.DataStore.GetRollCall(command.ChatID)
	if err != nil {
		return nil, err
	}
	if rollCall == nil {
		return &domain.Response{Text: "No roll call in progress", ChatID: command.ChatID}, nil
	}
	err = b.DataStore.SetQuiet(*rollCall, quiet)
	if err != nil {
		return nil, err
	}
	if quiet {
		return &domain.Response{ChatID: command.ChatID, Text: "Ok fine, I'll be quiet. 🤐"}, nil
	} else {
		err = b.DataStore.LoadRollCallResponses(rollCall)
		if err != nil {
			return nil, err
		}
		return &domain.Response{ChatID: command.ChatID, Text: "Sure. 😃\n" + responseList(rollCall)}, nil
	}
}

func (b *WhosInBot) handleWhosIn(command domain.Command) (*domain.Response, error) {
	rollCall, err := b.DataStore.GetRollCall(command.ChatID)
	if err != nil {
		return nil, err
	}
	if rollCall == nil {
		return &domain.Response{Text: "No roll call in progress", ChatID: command.ChatID}, nil
	}
	err = b.DataStore.LoadRollCallResponses(rollCall)
	if err != nil {
		return nil, err
	}
	return &domain.Response{ChatID: command.ChatID, Text: responseList(rollCall)}, nil
}

func (b *WhosInBot) handleResponse(command domain.Command, status string) (*domain.Response, error) {
	return b.setAttendanceFor(command, command.From.Name, status, command.ParamsString())
}

func (b *WhosInBot) handleResponseFor(command domain.Command, status string) (*domain.Response, error) {
	if len(command.FirstParam()) == 0 {
		return &domain.Response{ChatID: command.ChatID, Text: "Please provide the persons first name."}, nil
	}
	return b.setAttendanceFor(command, command.FirstParam(), status, command.ParamsStringExceptFirst())
}

func (b *WhosInBot) setAttendanceFor(command domain.Command, name string, status string, reason string) (*domain.Response, error) {
	rollCall, err := b.DataStore.GetRollCall(command.ChatID)
	if err != nil {
		return nil, err
	}
	if rollCall == nil {
		return &domain.Response{Text: "No roll call in progress", ChatID: command.ChatID}, nil
	}

	rollCallResponse := domain.NewRollCallResponse(command, name, status, reason)
	b.DataStore.SetResponse(rollCallResponse)

	err = b.DataStore.LoadRollCallResponses(rollCall)
	if err != nil {
		return nil, err
	}

	var response string
	if rollCall.Quiet {
		response = quietResponseList(rollCall, rollCallResponse)
	} else {
		response = responseList(rollCall)
	}
	return &domain.Response{ChatID: command.ChatID, Text: response}, nil
}

func responseList(rollCall *domain.RollCall) string {
	// DEBUG
	log.Printf("Response for roll call: %+v\n", rollCall)

	if len(rollCall.In) == 0 && len(rollCall.Out) == 0 && len(rollCall.Maybe) == 0 {
		return "No responses yet. 😢"
	}

	var text = ""

	if len(rollCall.Title) > 0 {
		text += rollCall.Title
	}

	if len(rollCall.In) > 0 && len(text) > 0 {
		text += "\n"
	}
	for index, response := range rollCall.In {
		text += fmt.Sprintf("%d. %v", index+1, response.Name)
		if len(response.Reason) > 0 {
			text += reasonWithParantheses(response.Reason)
		}
		if index+1 < len(rollCall.In) {
			text += "\n"
		}
	}

	text = appendResponses(text, rollCall.Out, "Out")
	text = appendResponses(text, rollCall.Maybe, "Maybe")

	return text
}

func quietResponseList(rollCall *domain.RollCall, rollCallResponse domain.RollCallResponse) string {
	var response = fmt.Sprintf("%v is %v!", rollCallResponse.Name, rollCallResponse.Status)
	return fmt.Sprintf("%v\nTotal: %d in, %d out, %d might come\n", response, len(rollCall.In), len(rollCall.Out), len(rollCall.Maybe))
}

func reasonWithParantheses(reason string) string {
	if strings.HasPrefix(reason, "(") {
		return " " + reason
	} else {
		return fmt.Sprintf(" (%v)", reason)
	}
}

func appendResponses(text string, responses []domain.RollCallResponse, status string) string {
	if len(responses) > 0 {
		if len(text) > 0 {
			text += "\n\n"
		}
		text += fmt.Sprintf("%v\n", status)
	}
	for index, response := range responses {
		text += fmt.Sprintf(" - %v", response.Name)
		if len(response.Reason) > 0 {
			text += fmt.Sprintf(" (%v)", response.Reason)
		}
		if index+1 < len(responses) {
			text += "\n"
		}
	}
	return text
}
