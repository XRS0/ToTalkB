package application

import (
	"log"
	"time"

	"event_manager/internal/domain"
)

type EventApplicationService struct {
	repository domain.EventRepository
	handlers   map[string]EventHandler
}

type EventHandler interface {
	Handle(event *domain.Event) error
}

func NewEventApplicationService(repository domain.EventRepository) *EventApplicationService {
	return &EventApplicationService{
		repository: repository,
		handlers:   make(map[string]EventHandler),
	}
}

func (s *EventApplicationService) RegisterHandler(eventType string, handler EventHandler) {
	s.handlers[eventType] = handler
}

func (s *EventApplicationService) ProcessEvent(event *domain.Event) error {
	event.Status = domain.StatusPending
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	if err := s.repository.Save(event); err != nil {
		return err
	}

	handler, exists := s.handlers[event.Type]
	if !exists {
		log.Printf("No handler registered for event type: %s", event.Type)
		return s.MarkAsFailed(event, "no handler registered")
	}

	if err := handler.Handle(event); err != nil {
		return s.MarkAsFailed(event, err.Error())
	}

	return s.MarkAsProcessed(event)
}

func (s *EventApplicationService) GetEventByID(id string) (*domain.Event, error) {
	return s.repository.FindByID(id)
}

func (s *EventApplicationService) MarkAsProcessed(event *domain.Event) error {
	event.Status = domain.StatusProcessed
	event.UpdatedAt = time.Now()
	return s.repository.Update(event)
}

func (s *EventApplicationService) MarkAsFailed(event *domain.Event, reason string) error {
	event.Status = domain.StatusFailed
	event.UpdatedAt = time.Now()
	return s.repository.Update(event)
}
