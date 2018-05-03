package core

import (
	"errors"
	"gopkg.in/telegram-bot-api.v4"
)

func HandleMessage(message *tgbotapi.Message) (string, error) {
	switch command := message.Command(); command {
	case "start_roll_call":
		return handleStart(message)
	case "end_roll_call":
		return handleEnd(message)
	case "in":
		return handleIn(message)
	case "out":
		return handleOut(message)
	default:
		return "", errors.New("Not a bot command")
	}
}

func handleStart(message *tgbotapi.Message) (string, error) {
	// TODO: Create new roll call
	return "Roll call started", nil
}

func handleEnd(message *tgbotapi.Message) (string, error) {
	// TODO: Delete roll call
	return "Roll call ended", nil
}

func handleIn(message *tgbotapi.Message) (string, error) {
	// TODO: mark sender as in
	return whosIn(message)
}

func handleOut(message *tgbotapi.Message) (string, error) {
	// TODO: mark sender as out
	return whosIn(message)
}

func whosIn(message *tgbotapi.Message) (string, error) {
	// TODO: output the roll call
	return "I don't know yet", nil
}
