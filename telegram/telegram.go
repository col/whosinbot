package telegram

import (
	"github.com/col/whosinbot/domain"
	"gopkg.in/telegram-bot-api.v4"
	"os"
	"fmt"
	"errors"
)

type Telegram struct {
	Token string
}

func NewTelegram(token string) (*Telegram) {
	return &Telegram{ Token: token }
}

func (t *Telegram) SendResponse(response *domain.Response) (error) {
	if response == nil || len(response.Text) == 0 {
		return nil
	}

	err := validateToken(t.Token)
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewMessage(response.ChatID, response.Text))
	if err != nil {
		return err
	}

	return nil
}

func validateToken(token string) (error) {
	validToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token != validToken {
		message := fmt.Sprintf("ERROR: Bot token doesn't match! Expected: %v Received: %v", validToken, token)
		return errors.New(message)
	}
	return nil
}