package interfaces

// RepositoryWriter определяет интерфейс для записи данных
type RepositoryWriter interface {
	Save(order *Order) error
	Init() error
}

// Notifier определяет интерфейс для отправки уведомлений
type Notifier interface {
	Notify(customer string, message string) error
}

type Order struct {
	ID       int
	Customer string
	Products string
	Total    float64
	Status   string
}
