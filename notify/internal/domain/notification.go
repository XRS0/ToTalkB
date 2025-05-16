package domain

import (
	"encoding/json"
	"time"
)

// Notification represents the core domain entity for notifications
type Notification struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Status    Status          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Status represents the current state of a notification
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	Save(notification *Notification) error
	FindByID(id string) (*Notification, error)
	Update(notification *Notification) error
}

// NotificationService defines the domain service interface
type NotificationService interface {
	ProcessNotification(notification *Notification) error
	MarkAsSent(notification *Notification) error
	MarkAsFailed(notification *Notification, reason string) error
}
