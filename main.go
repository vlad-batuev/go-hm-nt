package main

import (
	"database/sql"
	"fmt"
	"log"
	"refactor/interfaces"
	"refactor/notification"
	"refactor/repository"
)

// OrderService - основной сервис для работы с заказами
type OrderService struct {
	repo     interfaces.RepositoryWriter
	notifier interfaces.Notifier
}

func NewOrderService(repo interfaces.RepositoryWriter, notifier interfaces.Notifier) *OrderService {
	return &OrderService{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *OrderService) CreateOrder(customer string, products []string, total float64) error {
	// Создание объекта заказа
	order := &interfaces.Order{
		Customer: customer,
		Products: fmt.Sprintf("%v", products),
		Total:    total,
		Status:   "pending",
	}

	// Сохранение в БД
	if err := s.repo.Save(order); err != nil {
		return fmt.Errorf("ошибка сохранения заказа: %w", err)
	}

	// Отправка уведомления
	message := fmt.Sprintf("Ваш заказ на сумму %.2f создан", total)
	if err := s.notifier.Notify(customer, message); err != nil {
		return fmt.Errorf("ошибка отправки уведомления: %w", err)
	}

	return nil
}

func main() {
	// Инициализация базы данных
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание репозитория
	repo := repository.NewSQLiteOrderRepository(db)

	// Инициализация таблиц
	if err := repo.Init(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Пример 1: Использование EmailSender ===")
	emailNotifier := notification.NewEmailSender()
	emailService := NewOrderService(repo, emailNotifier)

	err = emailService.CreateOrder("Иван", []string{"apple", "banana"}, 10.5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Пример 2: Использование SMSSender ===")
	smsNotifier := notification.NewSMSSender()
	smsService := NewOrderService(repo, smsNotifier)

	err = smsService.CreateOrder("Мария", []string{"milk", "bread"}, 7.8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Пример 3: Легкое переключение между отправителями ===")

	// Можно динамически выбирать тип уведомления
	var notifier interfaces.Notifier

	// Выбираем email для VIP клиентов
	notifier = notification.NewEmailSender()
	vipService := NewOrderService(repo, notifier)
	vipService.CreateOrder("VIP Клиент", []string{"wine", "cheese"}, 50.0)

	// Выбираем SMS для обычных клиентов
	notifier = notification.NewSMSSender()
	regularService := NewOrderService(repo, notifier)
	regularService.CreateOrder("Обычный Клиент", []string{"water", "juice"}, 3.5)

	fmt.Println("\n✅ Все заказы успешно обработаны!")
}
