package main

import (
	"fmt"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/01-interfaces-and-polymorphism/main/notification"
)

func main() {
	fmt.Println("Hello, World!")
	emailNotifier := &notification.EmailNotifier{Sender: "viraj", Receiver: "unknown"}
	smsNotifier := &notification.SMSNotifier{PhoneNumber: "+1234567890"}
	slackNotifier := &notification.SlackNotifier{Channel: "#general", WebhookUrl: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}
	notifiers := []notification.Notifier{emailNotifier, smsNotifier, slackNotifier}

	notification.NotifyAll(notifiers, "System Alert: Server is down!")

}
