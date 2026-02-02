package repository

import (
	"context"
	"fmt"

	"L3.3/internal/models"
	"github.com/wb-go/wbf/dbpg"
)

// Repository интерфейс для хранилища
type Repository interface {
	// crud
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetUserById(ctx context.Context, userID int64) (*models.User, error)
	GetCommentById(ctx context.Context, commentID int64) (*models.Comment, error)
	DeleteUserById(ctx context.Context, userID int64) error
	DeleteCommentById(ctx context.Context, commentID int64) error

	//	методы для работы с деревьями
	// GetCommentsTree [] чтобы можно было отдавать сразу несколько корней
	GetCommentsTree(ctx context.Context, commentID int64) (*models.CommentNode, error)
	DeleteCommentsTree(ctx context.Context, commentID int64) error

	// Списки фильтры
	GetRootComments(ctx context.Context, limit, offset int64, sort string) ([]*models.Comment, error)
	SearchComments(ctx context.Context, query string, limit, offset int64) ([]*models.Comment, error)
}

// PostgresRepository стандартная реализация слоя repository
type PostgresRepository struct {
	db *dbpg.DB
}

// NewPostgresRepository создаёт новый репозиторий
func NewPostgresRepository(masterDSN string, opts *dbpg.Options) (*PostgresRepository, error) {
	if masterDSN == "" {
		return nil, fmt.Errorf("строка подключения к БД не найдена")
	}
	if opts == nil {
		return nil, fmt.Errorf("настройки пула подключений отсутствуют")
	}

	db, err := dbpg.New(masterDSN, nil, opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания подключения: %w", err)
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

// Close закрыть соединение
func (r *PostgresRepository) Close() error {
	if r.db != nil {
		return r.db.Master.Close()
	}
	return nil
}

// CreateUser создать пользователя
func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
	insert into users(username,email)
	values ($1,$2)
	returning id, created_at;
	`

	err := r.db.Master.QueryRowContext(ctx, query, user.Username, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка при добавлении пользователя: %w", err)
	}

	return user, nil
}

// CreateComment
func (r *PostgresRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	queryComment := `
		INSERT INTO comments (user_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at;
	`

	err = tx.QueryRowContext(ctx, queryComment, comment.UserID, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании комментария: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO comment_paths (ancestor_id, descendant_id, depth)
		 VALUES ($1, $1, 0);`,
		comment.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка вставки self-path: %w", err)
	}

	if comment.ParentID != nil {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO comment_paths (ancestor_id, descendant_id, depth)
			SELECT ancestor_id, $1, depth + 1
			FROM comment_paths
			WHERE descendant_id = $2;
		`, comment.ID, *comment.ParentID)
		if err != nil {
			return nil, fmt.Errorf("ошибка вставки путей родителя: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return comment, nil
}

// GetUserById возвращает пользователя
func (r *PostgresRepository) GetUserById(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1;
	`

	user := &models.User{}
	err := r.db.Master.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя с id=%d: %w", userID, err)
	}

	return user, nil
}

// GetCommentById возвращает комментарий
func (r *PostgresRepository) GetCommentById(ctx context.Context, commentID int64) (*models.Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.content, c.created_at, c.updated_at, cp.ancestor_id AS parent_id
		FROM comments c
		LEFT JOIN comment_paths cp ON cp.descendant_id = c.id AND cp.depth = 1
		WHERE c.id = $1;
	`

	comment := &models.Comment{}
	err := r.db.Master.QueryRowContext(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.ParentID,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении комментария с id=%d: %w", commentID, err)
	}

	return comment, nil
}

// DeleteUserById удаляет пользователя
func (r *PostgresRepository) DeleteUserById(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := r.db.Master.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя с id=%d: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества удалённых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с id=%d не найден", userID)
	}

	return nil
}

// DeleteCommentById удаляет комментарий
func (r *PostgresRepository) DeleteCommentById(ctx context.Context, commentID int64) error {
	tx, err := r.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Удаляем все комментарии, являющиеся потомками (включая сам commentID)
	queryDeleteComments := `
		DELETE FROM comments
		WHERE id IN (
			SELECT descendant_id FROM comment_paths WHERE ancestor_id = $1
		);
	`
	_, err = tx.ExecContext(ctx, queryDeleteComments, commentID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении комментариев дерева: %w", err)
	}

	// Удаляем пути из comment_paths (на всякий случай, если каскад не настроен)
	queryDeletePaths := `DELETE FROM comment_paths WHERE ancestor_id = $1 OR descendant_id = $1;`
	_, err = tx.ExecContext(ctx, queryDeletePaths, commentID)
	if err != nil {
		return fmt.Errorf("ошибка при очистке comment_paths: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при коммите удаления: %w", err)
	}

	return nil
}

// GetCommentsTree возвращает дерево комментариев
func (r *PostgresRepository) GetCommentsTree(ctx context.Context, commentID int64) (*models.CommentNode, error) {
	query := `
    SELECT c.id, c.user_id, c.content, c.created_at, cp.ancestor_id
    FROM comments c
    LEFT JOIN comment_paths cp ON c.id = cp.descendant_id AND cp.depth = 1
    WHERE c.id IN (
        SELECT descendant_id
        FROM comment_paths
        WHERE ancestor_id = $1
    )
    ORDER BY c.created_at;
    `

	rows, err := r.db.Master.QueryContext(ctx, query, commentID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении дерева комментариев для id=%v: %w", commentID, err)
	}
	defer rows.Close()

	nodes := make(map[int64]*models.CommentNode)
	parents := make(map[int64]*int64) // childID -> parentID
	var root *models.CommentNode

	for rows.Next() {
		node := &models.CommentNode{}
		var parentID *int64

		if err := rows.Scan(&node.ID, &node.UserID, &node.Content, &node.CreatedAt, &parentID); err != nil {
			return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
		}

		node.Children = []*models.CommentNode{}
		nodes[node.ID] = node
		if parentID != nil {
			pid := *parentID
			parents[node.ID] = &pid
		}
		if node.ID == commentID {
			root = node
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обходе строк: %w", err)
	}

	if root == nil {
		return nil, fmt.Errorf("комментарий с id=%d не найден", commentID)
	}

	// второй проход: связываем детей с родителями
	for childID, parentID := range parents {
		if parentID == nil {
			continue
		}
		parent := nodes[*parentID]
		child := nodes[childID]
		if parent != nil && child != nil {
			parent.Children = append(parent.Children, child)
		}
	}

	return root, nil
}

func (r *PostgresRepository) DeleteCommentsTree(ctx context.Context, commentID int64) error {
	query := `
	DELETE FROM comments
	WHERE id IN (
		SELECT descendant_id
		FROM comment_paths
		WHERE ancestor_id = $1
	);
	`

	_, err := r.db.Master.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении дерева комментариев: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetRootComments(ctx context.Context, limit, offset int64, sort string) ([]*models.Comment, error) {
	// защита от SQL-инъекций в параметре sort
	orderBy := "created_at"
	if sort == "desc" {
		orderBy = "created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.user_id, c.content, c.created_at, c.updated_at
		FROM comments c
		WHERE NOT EXISTS (
			SELECT 1 FROM comment_paths p
			WHERE p.descendant_id = c.id AND p.depth = 1
		)
		ORDER BY %s
		LIMIT $1 OFFSET $2;
	`, orderBy)

	rows, err := r.db.Master.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении корневых комментариев: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err = rows.Scan(&c.ID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении строки: %w", err)
		}
		comments = append(comments, &c)
	}

	return comments, nil
}

func (r *PostgresRepository) SearchComments(ctx context.Context, queryText string, limit, offset int64) ([]*models.Comment, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM comments
		WHERE to_tsvector('russian', content) @@ plainto_tsquery('russian', $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Master.QueryContext(ctx, query, queryText, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске комментариев: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err = rows.Scan(&c.ID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении строки: %w", err)
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
