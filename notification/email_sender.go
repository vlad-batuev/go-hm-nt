package notification

import "fmt"

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (e *EmailSender) Notify(customer string, message string) error {
	fmt.Printf("📧 Email отправлен клиенту %s: %s\n", customer, message)
	return nil
}
