package chat_bot

import (
	"encoding/json"
	"google.golang.org/api/chat/v1"
	"log"
	"strings"
	"whosinbot/domain"
)

func ParseDeprecatedEvent(requestBody []byte) (domain.Command, error) {
	deprecatedEvent := chat.DeprecatedEvent{}
	err := json.Unmarshal(requestBody, &deprecatedEvent)

	log.Printf("deprecatedEvent: " + deprecatedEvent.Message.Text)

	if err != nil {
		return domain.EmptyCommand(), err
	}
	// TODO: split this into a mapping function and write some tests around it
	threadName := deprecatedEvent.Message.Thread.Name
	argumentText := deprecatedEvent.Message.ArgumentText
	tokens := strings.Split(argumentText, " ")
	command := domain.Command{
		ChatID: threadName,
		Name: tokens[0],
		Params: []string{tokens[1]},
		From: domain.User{
			UserID: deprecatedEvent.Message.Sender.Name,
			Name:   deprecatedEvent.Message.Sender.DisplayName,
		},
	}
	return command, err
}
