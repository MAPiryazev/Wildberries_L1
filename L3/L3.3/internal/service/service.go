package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"L3.3/internal/config"
	customerrors "L3.3/internal/custom-errors"
	"L3.3/internal/models"
	"L3.3/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID int64, parentID *int64, content string) (*models.Comment, error)
	GetCommentTree(ctx context.Context, parentID *int64, page, limit int64, sort string) ([]*models.CommentNode, error)
	DeleteCommentTree(ctx context.Context, id int64, commentID int64) error
	SearchComments(ctx context.Context, query string, page, limit int64) ([]*models.Comment, error)
}

// commentService реализует интерфейс сервиса
type commentService struct {
	repo     repository.Repository
	maxDepth int64
}

// NewCommentService конструктор для сервиса
func NewCommentService(repo repository.Repository, config *config.APIConfig) *commentService {
	return &commentService{
		repo:     repo,
		maxDepth: config.MaxDepth,
	}
}

// CreateComment создаёт новый комментарий (с указанием родителя, если задан)
func (s *commentService) CreateComment(ctx context.Context, userID int64, parentID *int64, content string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("%w", customerrors.ErrCommentEmpty)
	}
	if len(content) > 1000 {
		return nil, fmt.Errorf("%w", customerrors.ErrCommentTooLong)
	}

	// Проверка существования пользователя
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя %d: %w", userID, err)
	}
	if user == nil {
		return nil, fmt.Errorf("%w: id=%d", customerrors.ErrNotFound, userID)
	}

	// Проверка родителя, если задан
	if parentID != nil {
		parent, err := s.repo.GetCommentById(ctx, *parentID)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				return nil, fmt.Errorf("%w: родительский комментарий id=%d", customerrors.ErrNotFound, *parentID)
			}
			return nil, fmt.Errorf("не удалось получить родительский комментарий: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("%w: родительский комментарий id=%d", customerrors.ErrNotFound, *parentID)
		}
	}

	comment := &models.Comment{
		UserID:    userID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	created, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать комментарий: %w", err)
	}

	return created, nil
}

// GetCommentTree получает комментарий и всех вложенных под ним
func (s *commentService) GetCommentTree(ctx context.Context, parentID *int64, page, limit int64, sort string) ([]*models.CommentNode, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Корневые комментарии
	if parentID == nil {
		roots, err := s.repo.GetRootComments(ctx, limit, offset, sort)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить корневые комментарии: %w", err)
		}

		var trees []*models.CommentNode
		for _, c := range roots {
			node, err := s.repo.GetCommentsTree(ctx, c.ID)
			if err != nil {
				return nil, fmt.Errorf("не удалось получить дерево комментария %d: %w", c.ID, err)
			}
			trees = append(trees, node)
		}
		return trees, nil
	}

	// Поддерево для указанного parentID
	tree, err := s.repo.GetCommentsTree(ctx, *parentID)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: id=%d", customerrors.ErrNotFound, *parentID)
		}
		return nil, fmt.Errorf("не удалось получить дерево комментариев: %w", err)
	}

	// Проверка глубины
	if depth := countDepthHelper(tree); int64(depth) > s.maxDepth {
		return nil, fmt.Errorf("%w: %d > %d", customerrors.ErrMaxDepthExceeded, depth, s.maxDepth)
	}

	return []*models.CommentNode{tree}, nil
}

// countDepth — вспомогательная функция для проверки глубины дерева
func countDepthHelper(node *models.CommentNode) int {
	if node == nil || len(node.Children) == 0 {
		return 1
	}
	maxChild := 0
	for _, c := range node.Children {
		if d := countDepthHelper(c); d > maxChild {
			maxChild = d
		}
	}
	return maxChild + 1
}

// DeleteCommentTree удаляет комментарий и все вложенные
func (s *commentService) DeleteCommentTree(ctx context.Context, userID, commentID int64) error {
	comment, err := s.repo.GetCommentById(ctx, commentID)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			return fmt.Errorf("%w: id=%d", customerrors.ErrNotFound, commentID)
		}
		return fmt.Errorf("не удалось получить комментарий %d: %w", commentID, err)
	}

	// Проверка прав на удаление (только автор может удалить)
	if comment.UserID != userID {
		return fmt.Errorf("%w: пользователь %d не может удалить комментарий %d", customerrors.ErrActionForbidden, userID, commentID)
	}

	if err := s.repo.DeleteCommentsTree(ctx, commentID); err != nil {
		return fmt.Errorf("не удалось удалить дерево комментариев %d: %w", commentID, err)
	}

	return nil
}

// SearchComments ищет комментарии по ключевому слову
func (s *commentService) SearchComments(ctx context.Context, query string, page, limit int64) ([]*models.Comment, error) {
	query = strings.TrimSpace(query)
	if len(query) < 3 {
		return nil, fmt.Errorf("%w", customerrors.ErrSearchQueryTooShort)
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	comments, err := s.repo.SearchComments(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить поиск комментариев: %w", err)
	}

	return comments, nil
}
