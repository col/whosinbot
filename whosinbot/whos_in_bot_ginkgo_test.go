package whosinbot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"whosinbot/domain"
	"fmt"
)

var _ = Describe("WhosInBot", func() {
	Describe("#handleAvailableCommands", func() {
		It("creates response with available commands besides itself", func() {
			command := domain.Command{
				ChatID: "123",
				Name:   "available_commands",
				Params: []string{},
				From: domain.User{
					UserID: "abc",
					Name:   "Sakamoto",
				},
			}

			formattedAvailableCommands := fmt.Sprint("start_roll_call\n" +
				"end_roll_call\n" +
				"set_title\n" +
				"in\n" +
				"out\n" +
				"maybe\n" +
				"set_in_for\n" +
				"set_out_for\n" +
				"set_maybe_for\n" +
				"whos_in\n" +
				"shh\n" +
				"louder\n")
			expectedResponse := &domain.Response{ChatID: command.ChatID,
				Text: formattedAvailableCommands}

			response, err := bot.handleAvailableCommands(command)
			Expect(response).To(BeEquivalentTo(expectedResponse))
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
