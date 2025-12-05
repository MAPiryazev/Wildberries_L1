package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	"github.com/wb-go/wbf/dbpg"
)

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

	log.Println("создано подключение к БД")

	return &PostgresRepository{
		db: db,
	}, nil
}

func (p *PostgresRepository) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if event == nil {
		return nil, fmt.Errorf("event не может быть nil при создании")
	}

	query := `
	insert into events (title, start_time, capacity)
	values ($1,$2,$3)
	returning id, title, start_time, capacity, created_at;
	`

	var created models.Event

	err := p.db.QueryRowContext(
		ctx,
		query,
		event.Title,
		event.StartTime,
		event.Capacity,
	).Scan(
		&created.ID,
		&created.Title,
		&created.StartTime,
		&created.Capacity,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка repository уровня при создании event: %w", err)
	}

	return &created, nil
}

func (p *PostgresRepository) GetEventByID(ctx context.Context, eventID int64) (*models.Event, error) {
	var event models.Event

	query := `
	SELECT id, title, start_time, capacity, created_at
	FROM events
	WHERE id = $1
	`

	err := p.db.QueryRowContext(ctx, query, eventID).Scan(
		&event.ID,
		&event.Title,
		&event.StartTime,
		&event.Capacity,
		&event.CreatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("запрошенный event с id %v не найден", eventID)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка repository при получении event по id: %w", err)
	}

	return &event, nil
}

// GetAllEvents возвращает все мероприятия
func (p *PostgresRepository) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	query := `
		SELECT id, title, start_time, capacity, created_at
		FROM events
		ORDER BY start_time
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка repository при получении всех events: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.StartTime, &e.Capacity, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании row: %w", err)
		}
		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время перебора rows: %w", err)
	}

	return events, nil
}

// CreateUser создаёт нового пользователя
func (p *PostgresRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user не может быть nil")
	}

	query := `
		INSERT INTO users (name, is_admin)
		VALUES ($1, $2)
		RETURNING id, name, is_admin
	`

	var created models.User
	err := p.db.QueryRowContext(ctx, query, user.Name, user.IsAdmin).Scan(
		&created.ID,
		&created.Name,
		&created.IsAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка repository при создании user: %w", err)
	}

	log.Printf("создан пользователь с ID = %v", created.ID)

	return &created, nil
}

// GetUserByID возвращает пользователя по id
func (p *PostgresRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, name, is_admin
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := p.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.IsAdmin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка repository при получении user по id: %w", err)
	}

	log.Printf("найден пользователь с ID = %v", user.ID)

	return &user, nil
}

// CreateBooking создаёт новую бронь, если есть свободные места
func (p *PostgresRepository) CreateBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	if booking == nil {
		return nil, fmt.Errorf("booking не может быть nil")
	}

	tx, err := p.db.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Считаем сколько уже подтверждено/забронировано мест
	var count int64
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM bookings 
		WHERE event_id = $1 AND status IN ('booked', 'confirmed')
	`, booking.EventID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка подсчёта занятых мест: %w", err)
	}

	// Получаем вместимость события
	var capacity int64
	err = tx.QueryRowContext(ctx, `SELECT capacity FROM events WHERE id = $1`, booking.EventID).Scan(&capacity)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения capacity события: %w", err)
	}

	if count >= capacity {
		return nil, fmt.Errorf("нет свободных мест для события %d", booking.EventID)
	}

	// Создаем бронь
	err = tx.QueryRowContext(ctx, `
		INSERT INTO bookings (event_id, user_id, status, expires_at)
		VALUES ($1, $2, 'booked', $3)
		RETURNING id, event_id, user_id, status, created_at, expires_at
	`, booking.EventID, booking.UserID, booking.ExpiresAt).Scan(
		&booking.ID,
		&booking.EventID,
		&booking.UserID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания брони: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return booking, nil
}

// GetBooking возвращает бронь по событию и пользователю
func (p *PostgresRepository) GetBooking(ctx context.Context, eventID, userID int64) (*models.Booking, error) {
	var b models.Booking
	err := p.db.QueryRowContext(ctx, `
		SELECT id, event_id, user_id, status, created_at, expires_at
		FROM bookings
		WHERE event_id = $1 AND user_id = $2
	`, eventID, userID).Scan(
		&b.ID,
		&b.EventID,
		&b.UserID,
		&b.Status,
		&b.CreatedAt,
		&b.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения брони: %w", err)
	}
	return &b, nil
}

// UpdateBookingStatus обновляет статус брони
func (p *PostgresRepository) UpdateBookingStatus(ctx context.Context, bookingID int64, status string) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE bookings
		SET status = $1
		WHERE id = $2
	`, status, bookingID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса брони: %w", err)
	}
	return nil
}

// CountActiveBookings возвращает количество активных бронирований (booked + confirmed) для события
func (p *PostgresRepository) CountActiveBookings(ctx context.Context, eventID int64) (int64, error) {
	var count int64
	err := p.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM bookings
		WHERE event_id = $1 AND status IN ('booked', 'confirmed')
	`, eventID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта активных бронирований: %w", err)
	}
	return count, nil
}

// FindExpiredBookings возвращает все брони со статусом 'booked', у которых истёк срок
func (p *PostgresRepository) FindExpiredBookings(ctx context.Context) ([]*models.Booking, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT id, event_id, user_id, status, created_at, expires_at
		FROM bookings
		WHERE status = 'booked' AND expires_at < NOW()
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения просроченных бронирований: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.EventID, &b.UserID, &b.Status, &b.CreatedAt, &b.ExpiresAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки брони: %w", err)
		}
		bookings = append(bookings, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк бронирований: %w", err)
	}

	return bookings, nil
}
