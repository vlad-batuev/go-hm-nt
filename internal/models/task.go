package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
