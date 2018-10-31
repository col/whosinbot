package hangout

import (
	"encoding/json"
	"google.golang.org/api/chat/v1"
	"log"
	"strings"
	"whosinbot/domain"
	"errors"
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
	arguments := strings.Fields(deprecatedEvent.Message.ArgumentText)
	if len(arguments) <= 0 {
		return domain.EmptyCommand(), errors.New("no argument provided")
	}
	command := domain.Command{
		ChatID: threadName,
		Name: arguments[0],
		Params: arguments[1:],
		From: domain.User{
			UserID: deprecatedEvent.Message.Sender.Name,
			Name:   deprecatedEvent.Message.Sender.DisplayName,
		},
	}
	return command, err
}
