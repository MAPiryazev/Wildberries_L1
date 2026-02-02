package models

import "time"

// CommentNode структура для использования в хендлерах
type CommentNode struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	ParentID  *int64         `json:"parent_id,omitempty"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	Children  []*CommentNode `json:"children,omitempty"`
}
