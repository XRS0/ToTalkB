package services

import (
	"context"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"
)

type EventQueueService struct {
	queueRepo domain.EventQueueRepository
}

func NewEventQueueService(repo domain.EventQueueRepository) *EventQueueService {
	return &EventQueueService{
		queueRepo: repo,
	}
}

func (s *EventQueueService) GetQueueStatus(ctx context.Context, eventID string) ([]*domain.EventQueue, error) {
	return s.queueRepo.GetByEventID(ctx, eventID)
}

func (s *EventQueueService) GetUserPosition(ctx context.Context, eventID, userID string) (int, error) {
	return s.queueRepo.GetUserPosition(ctx, eventID, userID)
}

func (s *EventQueueService) ProcessNext(ctx context.Context, eventID string) (*domain.EventQueue, error) {
	return s.queueRepo.ProcessNext(ctx, eventID)
}
