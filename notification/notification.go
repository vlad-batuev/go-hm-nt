package notification

import "refactor/interfaces"

// NotificationSender абстрактный интерфейс для отправки уведомлений
type NotificationSender interface {
	interfaces.Notifier
}
