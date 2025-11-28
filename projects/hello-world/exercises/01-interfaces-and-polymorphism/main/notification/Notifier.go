package notification

import "fmt"

type Notifier interface {
	Send(message string) error
	GetType() string
}

type EmailNotifier struct {
	Sender   string
	Receiver string
}

func (e *EmailNotifier) Send(message string) error {
	fmt.Println("Sending Email to ", e.Sender, " from ", e.Receiver, ": ", message)
	return nil
}

func (e *EmailNotifier) GetType() string {
	return "Email Notifier"
}

type SMSNotifier struct {
	PhoneNumber string
}

func (s *SMSNotifier) Send(message string) error {
	fmt.Println("Sending SMS to ", s.PhoneNumber, ": ", message)
	return nil
}

func (s *SMSNotifier) GetType() string {
	return "Sms Notifier"
}

type SlackNotifier struct {
	Channel    string
	WebhookUrl string
}

func (sl *SlackNotifier) Send(message string) error {
	fmt.Println("Sending Slack to ", sl.Channel, ": ", message)
	return nil
}

func (sl *SlackNotifier) GetType() string {
	return "Slack Notifier"
}

func NotifyAll(notifiers []Notifier, message string) {
	var valid bool = true
	for _, value := range notifiers {
		err := value.Send(message)
		if err != nil {
			valid = false
		}
	}
	if valid == true {
		fmt.Println("Notification sent successfully!")
	}
}
