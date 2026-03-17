package model

import "time"

type AuthUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "Completed"
)
