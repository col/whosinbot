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
	Describe("#ParseDeprecatedEvent", func() {
		When("successful", func() {
			It("parses request body to DeprecatedEvent", func() {
				deprecatedEvent := chat.DeprecatedEvent{
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
				requestBody, _ := json.Marshal(deprecatedEvent)
				command, _ := ParseDeprecatedEvent(requestBody)

				Expect(command.ChatID).To(Equal("thread1"))
				Expect(command.Name).To(Equal("cmdName"))
				Expect(command.Params).To(Equal([]string{"argument1", "argument2"}))
				Expect(command.From).To(BeEquivalentTo(domain.User{
					UserID: "users/1234567",
					Name:   "Ryuichi Sakamoto",
				}))
			})
		})

		When("failed", func() {
			Context("when there is no arguments provided", func() {
				deprecatedEvent := getEventWithNoArguments()

				It("returns empty command", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					command, _ := ParseDeprecatedEvent(requestBody)

					Expect(command).To(BeEquivalentTo(domain.EmptyCommand()))
				})

				It("returns error", func() {
					requestBody, _ := json.Marshal(deprecatedEvent)
					_, err := ParseDeprecatedEvent(requestBody)

					Expect(err).To(BeEquivalentTo( errors.New("no argument provided")))
				})
			})

			Context("when can't be unmarshaled to event", func() {
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
