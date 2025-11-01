package models

import "time"

// CommentNode структура для использования в хендлерах
type CommentNode struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	ParentID  *int           `json:"parent_id,omitempty"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	Children  []*CommentNode `json:"children,omitempty"`
}
