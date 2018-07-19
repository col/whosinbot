package http

import (
	"github.com/col/whosinbot/telegram"
	"github.com/col/whosinbot/whosinbot"
	"io/ioutil"
	"net/http"
	"os"
)

type WebhookHandler struct {
	WhosInBot *whosinbot.WhosInBot
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	command, err := telegram.ParseUpdate(body)
	if err != nil {
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	response, err := h.WhosInBot.HandleCommand(command)
	if err != nil {
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	api := telegram.NewTelegram(os.Getenv("TELEGRAM_BOT_TOKEN"))
	err = api.SendResponse(response)
	if err != nil {
		http.Error(w, "failed to send response", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
