package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetAll(ctx context.Context) ([]*domain.Event, error) {
	query := `SELECT id, type, source, payload, status, created_at, updated_at FROM events`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var event domain.Event
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&event.ID, &event.Type, &event.Source, &event.Payload, &event.Status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		event.CreatedAt = createdAt
		event.UpdatedAt = updatedAt
		events = append(events, &event)
	}
	return events, nil
}

func (r *EventRepository) List(ctx context.Context) ([]*domain.Event, error) {
	return r.GetAll(ctx)
}

func (r *EventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	query := `SELECT id, type, source, payload, status, created_at, updated_at FROM events WHERE id = $1`
	var event domain.Event
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(&event.ID, &event.Type, &event.Source, &event.Payload, &event.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	event.CreatedAt = createdAt
	event.UpdatedAt = updatedAt
	return &event, nil
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *EventRepository) Save(ctx context.Context, event *domain.Event) error {
	query := `INSERT INTO events (id, type, source, payload, status, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
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

func (r *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `UPDATE events 
			  SET type = $1, source = $2, payload = $3, status = $4, updated_at = $5 
			  WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		event.Type,
		event.Source,
		event.Payload,
		event.Status,
		event.UpdatedAt,
		event.ID,
	)
	return err
}
