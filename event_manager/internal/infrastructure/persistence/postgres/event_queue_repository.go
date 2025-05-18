package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"

	"github.com/google/uuid"
)

type EventQueueRepository struct {
	db *sql.DB
}

func NewEventQueueRepository(db *sql.DB) *EventQueueRepository {
	return &EventQueueRepository{db: db}
}

func (r *EventQueueRepository) Save(queue *domain.EventQueue) error {
	if queue.ID == "" {
		queue.ID = uuid.New().String()
	}
	queue.CreatedAt = time.Now()
	queue.UpdatedAt = time.Now()

	query := `
		INSERT INTO event_queues (id, event_id, user_id, status, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		queue.ID,
		queue.EventID,
		queue.UserID,
		queue.Status,
		queue.Position,
		queue.CreatedAt,
		queue.UpdatedAt,
	)

	return err
}

func (r *EventQueueRepository) FindByEventID(eventID string) ([]*domain.EventQueue, error) {
	query := `
		SELECT id, event_id, user_id, status, position, created_at, updated_at
		FROM event_queues
		WHERE event_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []*domain.EventQueue
	for rows.Next() {
		queue := &domain.EventQueue{}
		err := rows.Scan(
			&queue.ID,
			&queue.EventID,
			&queue.UserID,
			&queue.Status,
			&queue.Position,
			&queue.CreatedAt,
			&queue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		queues = append(queues, queue)
	}

	return queues, nil
}

func (r *EventQueueRepository) FindByUserID(userID string) ([]*domain.EventQueue, error) {
	query := `
		SELECT id, event_id, user_id, status, position, created_at, updated_at
		FROM event_queues
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []*domain.EventQueue
	for rows.Next() {
		queue := &domain.EventQueue{}
		err := rows.Scan(
			&queue.ID,
			&queue.EventID,
			&queue.UserID,
			&queue.Status,
			&queue.Position,
			&queue.CreatedAt,
			&queue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		queues = append(queues, queue)
	}

	return queues, nil
}

func (r *EventQueueRepository) Update(queue *domain.EventQueue) error {
	queue.UpdatedAt = time.Now()

	query := `
		UPDATE event_queues
		SET status = $1, position = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query,
		queue.Status,
		queue.Position,
		queue.UpdatedAt,
		queue.ID,
	)

	return err
}

func (r *EventQueueRepository) Delete(id string) error {
	query := `DELETE FROM event_queues WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventQueueRepository) GetNext(eventID string) (*domain.EventQueue, error) {
	query := `
		SELECT id, event_id, user_id, status, position, created_at, updated_at
		FROM event_queues
		WHERE event_id = $1 AND status = $2
		ORDER BY position ASC
		LIMIT 1
	`

	queue := &domain.EventQueue{}
	err := r.db.QueryRow(query, eventID, domain.QueueStatusWaiting).Scan(
		&queue.ID,
		&queue.EventID,
		&queue.UserID,
		&queue.Status,
		&queue.Position,
		&queue.CreatedAt,
		&queue.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return queue, nil
}

func (r *EventQueueRepository) GetPosition(eventID string, userID string) (int, error) {
	query := `
		SELECT position
		FROM event_queues
		WHERE event_id = $1 AND user_id = $2
	`

	var position int
	err := r.db.QueryRow(query, eventID, userID).Scan(&position)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("user not found in queue")
	}
	if err != nil {
		return 0, err
	}

	return position, nil
}
