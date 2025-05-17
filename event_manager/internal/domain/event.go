package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Payload   json.RawMessage `json:"payload"`
	Status    Status          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Status string

const (
	StatusPending   Status = "pending"
	StatusProcessed Status = "processed"
	StatusFailed    Status = "failed"
)

type EventRepository interface {
	Save(event *Event) error
	FindByID(id string) (*Event, error)
	Update(event *Event) error
}

type EventService interface {
	ProcessEvent(event *Event) error
	MarkAsProcessed(event *Event) error
	MarkAsFailed(event *Event, reason string) error
}
