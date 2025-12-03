package models

import "time"

type Event struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	Capacity  int64     `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}
