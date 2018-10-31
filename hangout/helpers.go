package hangout

import (
	"encoding/json"
	"google.golang.org/api/chat/v1"
	"log"
	"strings"
	"whosinbot/domain"
	"errors"
)

var (
	startRollCallAlias = "start"
	endRollCallAlias   = "end"
	whosInAlias        = "ls"
	setTitleAlias      = "title"
)

func ParseDeprecatedEvent(requestBody []byte) (domain.Command, error) {
	deprecatedEvent := chat.DeprecatedEvent{}
	err := json.Unmarshal(requestBody, &deprecatedEvent)
	if err != nil {
		return domain.EmptyCommand(), err
	}

	log.Printf("deprecatedEvent: " + deprecatedEvent.Message.Text)

	// TODO: split this into a mapping function and write some tests around it
	threadName := deprecatedEvent.Message.Thread.Name
	avaliableCommands := domain.Command{
		ChatID: threadName,
		Name:   "available_commands",
		Params: []string{},
		From: domain.User{
			UserID: deprecatedEvent.Message.Sender.Name,
			Name:   deprecatedEvent.Message.Sender.DisplayName,
		},
	}
	arguments := strings.Fields(deprecatedEvent.Message.ArgumentText)
	if len(arguments) <= 0 {
		return avaliableCommands, nil
	}

	name, err := parseCommandName(arguments[0])
	if err != nil {
		return avaliableCommands, nil
	}

	command := domain.Command{
		ChatID: threadName,
		Name:   name,
		Params: arguments[1:],
		From: domain.User{
			UserID: deprecatedEvent.Message.Sender.Name,
			Name:   deprecatedEvent.Message.Sender.DisplayName,
		},
	}
	return command, err
}

func contains(set []string, element string) bool {
	for _, oneEle := range set {
		if oneEle == element {
			return true
		}
	}
	return false
}

func parseCommandName(commandNameArgument string) (string, error) {
	if contains(domain.AllCommands(), commandNameArgument){
		return commandNameArgument, nil
	}

	alias := strings.ToLower(commandNameArgument)
	if alias == startRollCallAlias {
		return "start_roll_call", nil
	}

	if alias == endRollCallAlias {
		return "end_roll_call", nil
	}

	if alias == whosInAlias {
		return "whos_in", nil
	}

	if alias == setTitleAlias {
		return "set_title", nil
	}

	return "", errors.New("invalid")
}
