package telegram

import (
	"github.com/col/whosinbot/domain"
	"gopkg.in/telegram-bot-api.v4"
	"encoding/json"
	"strings"
)

func ParseUpdate(requestBody []byte) (domain.Command, error) {
	update := &tgbotapi.Update{}
	err := json.Unmarshal(requestBody, update)
	if err != nil {
		return domain.EmptyCommand(), err
	}
	// TODO: split this into a mapping function and write some tests around it
	command := domain.Command{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Command(),
		Params: strings.Fields(update.Message.CommandArguments()),
		From:   domain.User{
			UserID: int64(update.Message.From.ID),
			Name: update.Message.From.FirstName,
		},
	}
	return command, err
}