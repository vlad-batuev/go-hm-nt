package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
	"time"
)

type Handler struct {
	Store storage.Storage
}

func New(s storage.Storage) *Handler {
	return &Handler{Store: s}
}

func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		h.sendError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		h.sendError(w, "Title is required", http.StatusBadRequest)
		return
	}

	task.CreatedAt = time.Now().Format(time.RFC3339)
	task.Done = false

	createdTask, err := h.Store.Create(task)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, exists := h.Store.Get(id)
	if !exists {
		h.sendError(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedTask, err := h.Store.Update(id, task)
	if err != nil {
		if err.Error() == "task not found" {
			h.sendError(w, "Task not found", http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	err := h.Store.Delete(id)
	if err != nil {
		if err.Error() == "task not found" {
			h.sendError(w, "Task not found", http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
