package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/config"

	"github.com/XRS0/ToTalkB/notify/internal/domain"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg *config.DatabaseConfig) (*Repository, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Init() error {
	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infrastructure/persistence/postgres/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
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

func (r *Repository) GetScheduledNotifications(now time.Time) ([]*domain.Notification, error) {
	query := `
		SELECT id, type, payload, status, created_at, updated_at, scheduled_at
		FROM notifications
		WHERE scheduled_at <= $1 AND status = $2
	`

	rows, err := r.db.Query(query, now, domain.StatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var notification domain.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.Payload,
			&notification.Status,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&notification.ScheduledAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (r *Repository) CancelScheduledNotification(id string) error {
	query := `
		UPDATE notifications
		SET status = $1, updated_at = $2
		WHERE id = $3 AND status = $4 AND scheduled_at > $2
	`

	now := time.Now()
	result, err := r.db.Exec(query, domain.StatusCancelled, now, id, domain.StatusPending)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or cannot be cancelled")
	}

	return nil
}
