package customerrors

import "errors"

// ошибки репозитория
var (
	ErrNotFound      = errors.New("не найдено")
	ErrAlreadyExists = errors.New("уже существует")
	ErrInvalidInput  = errors.New("неправильный ввод")
)

// Ошибки сервисного слоя
var (
	ErrCommentTooLong      = errors.New("комментарий слишком длинный")
	ErrCommentEmpty        = errors.New("комментарий не может быть пустым")
	ErrCommentCannotDelete = errors.New("нельзя удалить комментарий: нет прав")
	ErrActionForbidden     = errors.New("действие запрещено")
	ErrSearchQueryTooShort = errors.New("поисковый запрос слишком короткий")
	ErrMaxDepthExceeded    = errors.New("достигнута максимальная глубина дерева комментариев")
)
