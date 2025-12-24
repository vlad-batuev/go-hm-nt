package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	"tasks-api/internal/http/middleware"
	"tasks-api/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	h := handlers.New(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", h.TasksCollection)
	mux.HandleFunc("/tasks/", h.TaskItem)
	mux.HandleFunc("/health", h.HealthCheck)

	handler := middleware.LoggingMiddleware(middleware.JSONMiddleware(mux))

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
