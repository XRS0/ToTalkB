package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/notification"
)

type EventApplicationService struct {
	repository domain.EventRepository
	notifier   *notification.Client
	handlers   map[string]EventHandler
	queues     map[string]*Queue
	mu         sync.RWMutex
}

type EventHandler interface {
	Handle(event *domain.Event) error
}

type Queue struct {
	ID        string
	EventID   string
	Users     []string
	Positions map[string]int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEventApplicationService(repository domain.EventRepository, notifier *notification.Client) *EventApplicationService {
	return &EventApplicationService{
		repository: repository,
		notifier:   notifier,
		handlers:   make(map[string]EventHandler),
		queues:     make(map[string]*Queue),
	}
}

func (s *EventApplicationService) RegisterHandler(eventType string, handler EventHandler) {
	s.handlers[eventType] = handler
}

func (s *EventApplicationService) ProcessEvent(ctx context.Context, eventType string, source string, payload []byte) (*domain.Event, error) {
	event := domain.NewEvent(eventType, source, payload)

	if err := s.repository.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to save event: %v", err)
	}

	// Send notification about the new event
	notificationPayload, err := json.Marshal(map[string]interface{}{
		"event_id": event.ID,
		"type":     event.Type,
		"source":   event.Source,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification payload: %v", err)
	}

	_, err = s.notifier.SendNotification(ctx, "event_created", notificationPayload)
	if err != nil {
		// Log the error but don't fail the event processing
		fmt.Printf("Failed to send notification: %v\n", err)
	}

	handler, exists := s.handlers[event.Type]
	if !exists {
		log.Printf("No handler registered for event type: %s", event.Type)
		return event, nil
	}

	if err := handler.Handle(event); err != nil {
		event.Status = "failed"
		if updateErr := s.repository.Update(ctx, event); updateErr != nil {
			log.Printf("Failed to update event status: %v", updateErr)
		}
		return event, err
	}

	event.Status = "processed"
	if err := s.repository.Update(ctx, event); err != nil {
		log.Printf("Failed to update event status: %v", err)
	}

	return event, nil
}

func (s *EventApplicationService) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *EventApplicationService) GetEventStatus(ctx context.Context, id string) (*domain.Event, error) {
	return s.repository.GetByID(ctx, id)
}

// Queue methods
func (s *EventApplicationService) JoinQueue(ctx context.Context, eventID string, userID string) (*domain.EventQueue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[eventID]
	if !exists {
		queue = &Queue{
			ID:        fmt.Sprintf("queue_%s", eventID),
			EventID:   eventID,
			Users:     make([]string, 0),
			Positions: make(map[string]int),
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.queues[eventID] = queue
	}

	// Check if user is already in queue
	if _, exists := queue.Positions[userID]; exists {
		return nil, fmt.Errorf("user already in queue")
	}

	// Add user to queue
	position := len(queue.Users) + 1
	queue.Users = append(queue.Users, userID)
	queue.Positions[userID] = position
	queue.UpdatedAt = time.Now()

	return &domain.EventQueue{
		ID:        queue.ID,
		EventID:   eventID,
		UserID:    userID,
		Status:    "waiting",
		Position:  position,
		CreatedAt: queue.CreatedAt,
		UpdatedAt: queue.UpdatedAt,
	}, nil
}

func (s *EventApplicationService) LeaveQueue(ctx context.Context, eventID string, userID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[eventID]
	if !exists {
		return false, fmt.Errorf("queue not found")
	}

	_, exists = queue.Positions[userID]
	if !exists {
		return false, fmt.Errorf("user not in queue")
	}

	// Remove user from queue
	delete(queue.Positions, userID)
	for i, id := range queue.Users {
		if id == userID {
			queue.Users = append(queue.Users[:i], queue.Users[i+1:]...)
			break
		}
	}

	// Update positions
	for i, id := range queue.Users {
		queue.Positions[id] = i + 1
	}

	queue.UpdatedAt = time.Now()
	return true, nil
}

func (s *EventApplicationService) GetQueueStatus(ctx context.Context, eventID string) ([]*domain.EventQueue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queue, exists := s.queues[eventID]
	if !exists {
		return nil, fmt.Errorf("queue not found")
	}

	queues := make([]*domain.EventQueue, len(queue.Users))
	for i, userID := range queue.Users {
		queues[i] = &domain.EventQueue{
			ID:        queue.ID,
			EventID:   eventID,
			UserID:    userID,
			Status:    "waiting",
			Position:  queue.Positions[userID],
			CreatedAt: queue.CreatedAt,
			UpdatedAt: queue.UpdatedAt,
		}
	}

	return queues, nil
}

func (s *EventApplicationService) GetUserPosition(ctx context.Context, eventID string, userID string) (int32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queue, exists := s.queues[eventID]
	if !exists {
		return 0, fmt.Errorf("queue not found")
	}

	position, exists := queue.Positions[userID]
	if !exists {
		return 0, fmt.Errorf("user not in queue")
	}

	return int32(position), nil
}

func (s *EventApplicationService) ProcessNext(ctx context.Context, eventID string) (*domain.EventQueue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[eventID]
	if !exists {
		return nil, fmt.Errorf("queue not found")
	}

	if len(queue.Users) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Get next user
	userID := queue.Users[0]
	queue.Users = queue.Users[1:]
	delete(queue.Positions, userID)

	// Update positions
	for i, id := range queue.Users {
		queue.Positions[id] = i + 1
	}

	queue.UpdatedAt = time.Now()

	return &domain.EventQueue{
		ID:        queue.ID,
		EventID:   eventID,
		UserID:    userID,
		Status:    "completed",
		Position:  0,
		CreatedAt: queue.CreatedAt,
		UpdatedAt: queue.UpdatedAt,
	}, nil
}

func (s *EventApplicationService) CloseQueue(ctx context.Context, eventID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[eventID]
	if !exists {
		return false, fmt.Errorf("queue not found")
	}

	queue.Status = "closed"
	queue.UpdatedAt = time.Now()
	return true, nil
}
