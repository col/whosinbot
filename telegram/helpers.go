package telegram

import (
	"encoding/json"
	"whosinbot/domain"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"strconv"
)

func ParseUpdate(requestBody []byte) (domain.Command, error) {
	update := &tgbotapi.Update{}
	err := json.Unmarshal(requestBody, update)
	if err != nil {
		return domain.EmptyCommand(), err
	}
	// TODO: split this into a mapping function and write some tests around it
	command := domain.Command{
		ChatID: strconv.Itoa(int(update.Message.Chat.ID)),
		Name:   update.Message.Command(),
		Params: strings.Fields(update.Message.CommandArguments()),
		From: domain.User{
			UserID: strconv.Itoa(update.Message.From.ID),
			Name:   update.Message.From.FirstName,
		},
	}
	return command, err
}
