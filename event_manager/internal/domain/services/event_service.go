package services

import (
	"context"
	"fmt"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"
	"github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/notification"
)

type EventService struct {
	eventRepo domain.EventRepository
	notifier  *notification.Client
}

func NewEventService(repo domain.EventRepository, notifier *notification.Client) *EventService {
	return &EventService{
		eventRepo: repo,
		notifier:  notifier,
	}
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]*gen.Event, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	result := make([]*gen.Event, 0, len(events))
	for _, event := range events {
		result = append(result, &gen.Event{
			Id:        event.ID,
			Type:      event.Type,
			Source:    event.Source,
			Payload:   event.Payload,
			Status:    event.Status,
			CreatedAt: event.CreatedAt.Format(time.RFC3339),
			UpdatedAt: event.UpdatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
