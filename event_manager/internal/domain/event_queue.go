package domain

import (
	"context"
	"time"
)

// EventQueue представляет запись в очереди событий
type EventQueue struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventQueueStatus представляет возможные статусы записи в очереди
type EventQueueStatus string

const (
	QueueStatusWaiting   EventQueueStatus = "waiting"
	QueueStatusActive    EventQueueStatus = "active"
	QueueStatusCompleted EventQueueStatus = "completed"
	QueueStatusCancelled EventQueueStatus = "cancelled"
)

// EventQueueRepository определяет методы для работы с очередью событий
type EventQueueRepository interface {
	// Save сохраняет запись в очереди
	Save(queue *EventQueue) error
	// FindByEventID находит все записи в очереди для конкретного события
	FindByEventID(eventID string) ([]*EventQueue, error)
	// FindByUserID находит все записи в очереди для конкретного пользователя
	FindByUserID(userID string) ([]*EventQueue, error)
	// Update обновляет запись в очереди
	Update(queue *EventQueue) error
	// Delete удаляет запись из очереди
	Delete(id string) error
	// GetNext получает следующую запись в очереди для события
	GetNext(eventID string) (*EventQueue, error)
	// GetPosition получает позицию пользователя в очереди
	GetPosition(eventID string, userID string) (int, error)
	// GetByEventID находит все записи в очереди для конкретного события
	GetByEventID(ctx context.Context, eventID string) ([]*EventQueue, error)
	// GetUserPosition получает позицию пользователя в очереди
	GetUserPosition(ctx context.Context, eventID, userID string) (int, error)
	// ProcessNext обрабатывает следующую запись в очереди
	ProcessNext(ctx context.Context, eventID string) (*EventQueue, error)
}

// EventQueueService определяет бизнес-логику для работы с очередью событий
type EventQueueService interface {
	// JoinQueue добавляет пользователя в очередь события
	JoinQueue(eventID string, userID string) error
	// LeaveQueue удаляет пользователя из очереди события
	LeaveQueue(eventID string, userID string) error
	// GetQueueStatus получает статус очереди для события
	GetQueueStatus(eventID string) ([]*EventQueue, error)
	// GetUserPosition получает позицию пользователя в очереди
	GetUserPosition(eventID string, userID string) (int, error)
	// ProcessNext обрабатывает следующую запись в очереди
	ProcessNext(eventID string) error
	// CloseQueue закрывает набор в очередь для события
	CloseQueue(eventID string) error
}
