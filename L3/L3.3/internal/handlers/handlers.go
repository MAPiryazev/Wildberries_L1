package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"

	customerrors "L3.3/internal/custom-errors"
	"L3.3/internal/service"
)

// CommentHandler интерфейс для HTTP хендлеров комментариев
type CommentHandler interface {
	RegisterRoutes(r *ginext.RouterGroup)
	CreateComment(c *ginext.Context)
	GetComments(c *ginext.Context)
	DeleteComment(c *ginext.Context)
	SearchComments(c *ginext.Context)
}

// commentHandler реализация хендлера
type commentHandler struct {
	service service.CommentService
}

// NewCommentHandler конструктор
func NewCommentHandler(s service.CommentService) CommentHandler {
	return &commentHandler{service: s}
}

// RegisterRoutes подключает маршруты к RouterGroup
func (h *commentHandler) RegisterRoutes(r *ginext.RouterGroup) {
	r.POST("/comments", h.CreateComment)
	r.GET("/comments", h.GetComments)
	r.DELETE("/comments/:id", h.DeleteComment)
	r.GET("/comments/search", h.SearchComments)
}

// CreateComment обрабатывает POST /comments
func (h *commentHandler) CreateComment(c *ginext.Context) {
	var req struct {
		UserID   int64  `json:"user_id"`
		ParentID *int64 `json:"parent_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "недопустимое тело запроса"})
		return
	}

	comment, err := h.service.CreateComment(c, req.UserID, req.ParentID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrCommentEmpty),
			errors.Is(err, customerrors.ErrCommentTooLong),
			errors.Is(err, customerrors.ErrSearchQueryTooShort):
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		case errors.Is(err, customerrors.ErrNotFound):
			c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments обрабатывает GET /comments?parent={id}&page=&limit=&sort=
func (h *commentHandler) GetComments(c *ginext.Context) {
	parentIDStr := c.Query("parent")
	var parentID *int64
	if parentIDStr != "" {
		id, err := strconv.ParseInt(parentIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "неверный id родителя"})
			return
		}
		parentID = &id
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	sort := c.DefaultQuery("sort", "asc")

	comments, err := h.service.GetCommentTree(c, parentID, page, limit, sort)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrNotFound):
			c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		case errors.Is(err, customerrors.ErrMaxDepthExceeded):
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment обрабатывает DELETE /comments/:id
func (h *commentHandler) DeleteComment(c *ginext.Context) {
	idStr := c.Param("id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "неверный id комментария"})
		return
	}

	// В реальном приложении userID берётся из авторизации
	// Здесь для примера — query-параметр user_id
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "user_id не может быть пустым"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "неверный id пользователя"})
		return
	}

	err = h.service.DeleteCommentTree(c.Request.Context(), userID, commentID)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrNotFound):
			c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		case errors.Is(err, customerrors.ErrActionForbidden):
			c.JSON(http.StatusForbidden, ginext.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchComments обрабатывает GET /comments/search?query=&page=&limit=
func (h *commentHandler) SearchComments(c *ginext.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "запрос не может быть пустым"})
		return
	}

	if len(query) < 3 {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "запрос слишком короткий, минимум 3 символа"})
		return
	}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	comments, err := h.service.SearchComments(c.Request.Context(), query, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrSearchQueryTooShort):
			c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comments)
}
