package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/domain"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Save(notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, payload, status, created_at, updated_at, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	payload, err := json.Marshal(notification.Payload)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query,
		notification.ID,
		notification.UserID,
		notification.Type,
		payload,
		notification.Status,
		notification.CreatedAt,
		notification.UpdatedAt,
		notification.ScheduledAt,
	)

	return err
}

func (r *NotificationRepository) FindByID(id string) (*domain.Notification, error) {
	query := `
		SELECT id, user_id, type, payload, status, created_at, updated_at, scheduled_at
		FROM notifications
		WHERE id = $1
	`

	var notification domain.Notification
	var payload []byte
	var scheduledAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&payload,
		&notification.Status,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&scheduledAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	notification.Payload = payload
	if scheduledAt.Valid {
		notification.ScheduledAt = &scheduledAt.Time
	}

	return &notification, nil
}

func (r *NotificationRepository) Update(notification *domain.Notification) error {
	query := `
		UPDATE notifications
		SET user_id = $1, type = $2, payload = $3, status = $4, updated_at = $5, scheduled_at = $6
		WHERE id = $7
	`

	payload, err := json.Marshal(notification.Payload)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query,
		notification.UserID,
		notification.Type,
		payload,
		notification.Status,
		notification.UpdatedAt,
		notification.ScheduledAt,
		notification.ID,
	)

	return err
}

func (r *NotificationRepository) GetScheduledNotifications(now time.Time) ([]*domain.Notification, error) {
	query := `
		SELECT id, user_id, type, payload, status, created_at, updated_at, scheduled_at
		FROM notifications
		WHERE status = $1 AND scheduled_at <= $2
	`

	rows, err := r.db.Query(query, domain.StatusPending, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var notification domain.Notification
		var payload []byte
		var scheduledAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&payload,
			&notification.Status,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&scheduledAt,
		)
		if err != nil {
			return nil, err
		}

		notification.Payload = payload
		if scheduledAt.Valid {
			notification.ScheduledAt = &scheduledAt.Time
		}

		notifications = append(notifications, &notification)
	}

	return notifications, rows.Err()
}

func (r *NotificationRepository) CancelScheduledNotification(id string) error {
	query := `
		UPDATE notifications
		SET status = $1, updated_at = $2
		WHERE id = $3 AND status = $4
	`

	_, err := r.db.Exec(query,
		domain.StatusCancelled,
		time.Now(),
		id,
		domain.StatusPending,
	)

	return err
}
