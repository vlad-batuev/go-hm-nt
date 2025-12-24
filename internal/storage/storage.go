package storage

import "tasks-api/internal/models"

type Storage interface {
	List() []models.Task
	Create(task models.Task) (models.Task, error) // Добавлен параметр name
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error) // Добавлен параметр name
	Delete(id int) error
}
