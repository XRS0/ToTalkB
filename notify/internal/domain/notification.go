package domain

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Status    Status          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type NotificationRepository interface {
	Save(notification *Notification) error
	FindByID(id string) (*Notification, error)
	Update(notification *Notification) error
}

type NotificationService interface {
	ProcessNotification(notification *Notification) error
	MarkAsSent(notification *Notification) error
	MarkAsFailed(notification *Notification, reason string) error
}
