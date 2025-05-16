package postgres

import (
	"database/sql"
	"time"

	"notify/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Save(notification *domain.Notification) error {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	query := `
		INSERT INTO notifications (id, type, payload, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		notification.ID,
		notification.Type,
		notification.Payload,
		notification.Status,
		notification.CreatedAt,
		notification.UpdatedAt,
	)

	return err
}

func (r *Repository) FindByID(id string) (*domain.Notification, error) {
	query := `
		SELECT id, type, payload, status, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	var notification domain.Notification
	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.Type,
		&notification.Payload,
		&notification.Status,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *Repository) Update(notification *domain.Notification) error {
	notification.UpdatedAt = time.Now()

	query := `
		UPDATE notifications
		SET type = $1, payload = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query,
		notification.Type,
		notification.Payload,
		notification.Status,
		notification.UpdatedAt,
		notification.ID,
	)

	return err
}
