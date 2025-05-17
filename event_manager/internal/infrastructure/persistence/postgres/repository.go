package postgres

import (
	"database/sql"
	"fmt"

	"event_manager/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {
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

func (r *Repository) Save(event *domain.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	query := `
		INSERT INTO events (id, type, source, payload, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		event.ID,
		event.Type,
		event.Source,
		event.Payload,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	)

	return err
}

func (r *Repository) FindByID(id string) (*domain.Event, error) {
	query := `
		SELECT id, type, source, payload, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	event := &domain.Event{}
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Type,
		&event.Source,
		&event.Payload,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *Repository) Update(event *domain.Event) error {
	query := `
		UPDATE events
		SET type = $1, source = $2, payload = $3, status = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(query,
		event.Type,
		event.Source,
		event.Payload,
		event.Status,
		event.UpdatedAt,
		event.ID,
	)

	return err
}
