package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"calendar/internal/config"
	"calendar/internal/models"

	_ "github.com/lib/pq"
)

// PostgresStorage реализует хранилище в соответствии с интерфейсом storage
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage конструктор для PostgresStorage
func NewPostgresStorage(cfg *config.DBPSQLConfig) (*PostgresStorage, error) {
	dbCredentials := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dbCredentials)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", "ошибка создания подключения", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifeTime) * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", "не удалось подключиться к БД", err)
	}

	return &PostgresStorage{db: db}, nil
}

// CreateEvent функция для создания записи в календаре
func (p *PostgresStorage) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
		INSERT INTO events (user_id, date, title, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	row := p.db.QueryRowContext(ctx, query, event.UserID, event.Date, event.Title)
	err := row.Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать событие: %w", err)
	}

	return event, nil
}

// GetEvent функция для получения записи в календаре
func (p *PostgresStorage) GetEvent(ctx context.Context, eventID int) (*models.Event, error) {
	query := `SELECT id, user_id, date, title, created_at, updated_at FROM events WHERE id=$1`
	event := &models.Event{}
	err := p.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID, &event.UserID, &event.Date, &event.Title, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("событие не найдено")
		}
		return nil, fmt.Errorf("не удалось получить событие: %w", err)
	}

	return event, nil
}

// UpdateEvent функция для обновления записи в календаре
func (p *PostgresStorage) UpdateEvent(ctx context.Context, event *models.Event) error {
	query := `
		UPDATE events
		SET user_id=$1, date=$2, title=$3, updated_at=NOW()
		WHERE id=$4
	`
	res, err := p.db.ExecContext(ctx, query, event.UserID, event.Date, event.Title, event.ID)
	if err != nil {
		return fmt.Errorf("не удалось обновить событие: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновлённых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("событие не найдено")
	}

	return nil
}

// DeleteEvent функция для удаления записи из календаря
func (p *PostgresStorage) DeleteEvent(ctx context.Context, eventID int) error {
	query := `DELETE FROM events WHERE id=$1`
	res, err := p.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("не удалось удалить событие: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("событие не найдено")
	}

	return nil
}

func (p *PostgresStorage) eventsByRange(ctx context.Context, userID int, start, end time.Time) ([]*models.Event, error) {
	query := `SELECT id, user_id, date, title, created_at, updated_at 
			  FROM events
			  WHERE user_id=$1 AND date >= $2 AND date <= $3
			  ORDER BY date`
	rows, err := p.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить события: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		ev := &models.Event{}
		if err := rows.Scan(&ev.ID, &ev.UserID, &ev.Date, &ev.Title, &ev.CreatedAt, &ev.UpdatedAt); err != nil {
			return nil, fmt.Errorf("не удалось считать событие: %w", err)
		}
		events = append(events, ev)
	}

	return events, nil
}

// EventsForDay функция для получения всех событий за ближайший день
func (p *PostgresStorage) EventsForDay(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты")
	}
	start := date
	end := date
	return p.eventsByRange(ctx, userID, start, end)
}

// EventsForWeek функция для получения всех событий за ближайшую неделю
func (p *PostgresStorage) EventsForWeek(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты")
	}
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := date.AddDate(0, 0, -weekday+1)
	end := start.AddDate(0, 0, 6)
	return p.eventsByRange(ctx, userID, start, end)
}

// EventsForMonth функция для получения всех событий за ближайший месяц
func (p *PostgresStorage) EventsForMonth(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты")
	}
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, -1)
	return p.eventsByRange(ctx, userID, start, end)
}
