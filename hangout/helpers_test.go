package hangout

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/chat/v1"
	"encoding/json"
	"whosinbot/domain"
	"errors"
)

var _ = Describe("helper", func() {
	Describe("#parseCommandName", func() {
		When("command name is in the commands set", func() {
			It("returns the original name", func() {
				commandNmae, err := parseCommandName("start_roll_call")
				Expect(commandNmae).To(Equal("start_roll_call"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("end_roll_call")
				Expect(commandNmae).To(Equal("end_roll_call"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("set_title")
				Expect(commandNmae).To(Equal("set_title"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("in")
				Expect(commandNmae).To(Equal("in"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("out")
				Expect(commandNmae).To(Equal("out"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("maybe")
				Expect(commandNmae).To(Equal("maybe"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("set_in_for")
				Expect(commandNmae).To(Equal("set_in_for"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("set_out_for")
				Expect(commandNmae).To(Equal("set_out_for"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("set_maybe_for")
				Expect(commandNmae).To(Equal("set_maybe_for"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("whos_in")
				Expect(commandNmae).To(Equal("whos_in"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("shh")
				Expect(commandNmae).To(Equal("shh"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("louder")
				Expect(commandNmae).To(Equal("louder"))
				Expect(err).NotTo(HaveOccurred())

			})
		})

		When("command argument is an alias of a command", func() {
			It("returns original command name", func() {
				commandNmae, err := parseCommandName("start")
				Expect(commandNmae).To(Equal("start_roll_call"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("end")
				Expect(commandNmae).To(Equal("end_roll_call"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("ls")
				Expect(commandNmae).To(Equal("whos_in"))
				Expect(err).NotTo(HaveOccurred())

				commandNmae, err = parseCommandName("title")
				Expect(commandNmae).To(Equal("set_title"))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("command argument not found", func() {
			It("returns \"\" and error", func() {
				commandNmae, err := parseCommandName("blah")
				Expect(commandNmae).To(Equal(""))
				Expect(err).To(BeEquivalentTo(errors.New("invalid")))
			})
		})
	})

	Describe("#ParseDeprecatedEvent", func() {
		When("successful", func() {
			deprecatedEvent := chat.DeprecatedEvent{
				Action:                    nil,
				ConfigCompleteRedirectUrl: "",
				EventTime:                 "",
				Message: &chat.Message{
					ArgumentText: "   in argument1   argument2",
					Sender: &chat.User{
						DisplayName: "Ryuichi Sakamoto",
						Name:        "users/1234567",
					},
					Thread: &chat.Thread{
						Name: "thread1",
					},
				},
			}

			It("parses request body to DeprecatedEvent", func() {
				requestBody, _ := json.Marshal(deprecatedEvent)
				command, _ := ParseDeprecatedEvent(requestBody)

				Expect(command.ChatID).To(Equal("thread1"))
				Expect(command.Name).To(Equal("in"))
				Expect(command.Params).To(Equal([]string{"argument1", "argument2"}))
				Expect(command.From).To(BeEquivalentTo(domain.User{
					UserID: "users/1234567",
					Name:   "Ryuichi Sakamoto",
				}))
			})

			Context("parses alias to valid command name", func() {
				It("replaces 'start' as start_roll_call command", func() {
					deprecatedEvent := mockEvent()
					deprecatedEvent.Message.ArgumentText = "start foo bar"

					requestBody, _ := json.Marshal(deprecatedEvent)
					command, _ := ParseDeprecatedEvent(requestBody)

					Expect(command.Name).To(Equal("start_roll_call"))
					Expect(command.Params).To(Equal([]string{"foo", "bar"}))
				})
			})

			When("there is no arguments provided", func() {
				deprecatedEvent := getEventWithNoArguments()

				It("returns list all command", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					command, _ := ParseDeprecatedEvent(requestBody)

					Expect(command.ChatID).To(Equal("thread1"))
					Expect(command.Name).To(Equal("available_commands"))
					Expect(command.From).To(BeEquivalentTo(domain.User{
						UserID: "users/1234567",
						Name:   "Ryuichi Sakamoto",
					}))
				})

				It("don't returns error", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					_, err := ParseDeprecatedEvent(requestBody)

					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("fails to get the command name", func() {
				deprecatedEvent := mockEvent()
				deprecatedEvent.Message.ArgumentText = "blah"

				It("returns list all command", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					command, _ := ParseDeprecatedEvent(requestBody)

					Expect(command.Name).To(Equal("available_commands"))
				})

				It("don't returns error", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					_, err := ParseDeprecatedEvent(requestBody)

					Expect(err).NotTo(HaveOccurred())
				})

			})
		})

		When("failed", func() {
			When("can't be unmarshaled to event", func() {
				requestBody, _ := json.Marshal([]byte("something"))

				It("returns empty command", func() {
					command, _ := ParseDeprecatedEvent(requestBody)

					Expect(command).To(BeEquivalentTo(domain.EmptyCommand()))
				})

				It("returns error", func() {
					_, err := ParseDeprecatedEvent(requestBody)

					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})

func mockEvent() chat.DeprecatedEvent {
	return chat.DeprecatedEvent{
		Action:                    nil,
		ConfigCompleteRedirectUrl: "",
		EventTime:                 "",
		Message: &chat.Message{
			ArgumentText: "   cmdName argument1   argument2",
			Sender: &chat.User{
				DisplayName: "Ryuichi Sakamoto",
				Name:        "users/1234567",
			},
			Thread: &chat.Thread{
				Name: "thread1",
			},
		},
	}
}

func getEventWithNoArguments() chat.DeprecatedEvent {
	return chat.DeprecatedEvent{
		Action:                    nil,
		ConfigCompleteRedirectUrl: "",
		EventTime:                 "",
		Message: &chat.Message{
			Sender: &chat.User{
				DisplayName: "Ryuichi Sakamoto",
				Name:        "users/1234567",
			},
			Thread: &chat.Thread{
				Name: "thread1",
			},
		},
	}
}
