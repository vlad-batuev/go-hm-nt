package storage

import (
	"errors"
	"sync"
	"tasks-api/internal/models"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *MemoryStorage) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.Title == "" {
		return models.Task{}, errors.New("title is required")
	}

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStorage) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *MemoryStorage) Update(id int, updatedTask models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return models.Task{}, errors.New("task not found")
	}

	if updatedTask.Title == "" {
		return models.Task{}, errors.New("title is required")
	}

	updatedTask.ID = id
	if updatedTask.CreatedAt == "" {
		updatedTask.CreatedAt = s.tasks[id].CreatedAt
	}
	s.tasks[id] = updatedTask
	return updatedTask, nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return errors.New("task not found")
	}

	delete(s.tasks, id)
	return nil
}
