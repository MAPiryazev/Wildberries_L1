package entity

import "time"

// UserID - псевдоним для user ключа.
type UserID int

// EventID - псевдоним для event ключа
type EventID int

// Event описывает событие календаря.
type Event struct {
	EventID  EventID    `json:"event_id"`
	UserID   UserID     `json:"user_id"`
	Date     time.Time  `json:"date"`
	Title    string     `json:"title"`
	RemindAt *time.Time `json:"remind_at,omitempty"`
	Archived bool       `json:"-"`
}
