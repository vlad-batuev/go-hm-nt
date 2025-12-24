package notification

import "fmt"

type SMSSender struct{}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

func (s *SMSSender) Notify(customer string, message string) error {
	fmt.Printf("📱 SMS отправлено клиенту %s: %s\n", customer, message)
	return nil
}
