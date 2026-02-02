package models

import "time"

// Event представляет собой сущность "событие" в календаре
type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      time.Time `json:"date"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
