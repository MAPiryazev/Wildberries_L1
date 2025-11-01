package repository

import (
	"context"

	"L3.3/internal/models"
)

// Repository интерфейс для хранилища
type Repository interface {
	// crud
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetUserById(ctx context.Context, UserID int) (*models.User, error)
	GetCommentById(ctx context.Context, CommentID int) (*models.Comment, error)
	DeleteUserById(ctx context.Context, UserID int) error
	DeleteCommentById(ctx context.Context, CommentID int) error

	//	методы для работы с деревьями
	// GetCommentsTree [] чтобы можно было отдавать сразу несколько корней
	GetCommentsTree(ctx context.Context, CommentID int) ([]*models.CommentNode, error)
	DeleteCommentsTree(ctx context.Context, CommentID int) error

	// Списки фильтры
	GetRootComments(ctx context.Context, limit, offset int, sort string) ([]*models.Comment, error)
	SearchComments(ctx context.Context, query string, limit, offset int) ([]*models.Comment, error)
}
