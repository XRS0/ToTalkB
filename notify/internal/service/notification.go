package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/domain"
)

type NotificationService struct {
	handlers   map[string]NotificationHandler
	repository domain.NotificationRepository
	mu         sync.RWMutex
	stopCh     chan struct{}
}

type NotificationHandler interface {
	Handle(notification *domain.Notification) error
}

func NewNotificationService(repository domain.NotificationRepository) *NotificationService {
	return &NotificationService{
		handlers:   make(map[string]NotificationHandler),
		repository: repository,
		stopCh:     make(chan struct{}),
	}
}

func (s *NotificationService) RegisterHandler(notificationType string, handler NotificationHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[notificationType] = handler
}

func (s *NotificationService) ProcessNotification(notification *domain.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repository.Save(notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	handler, exists := s.handlers[notification.Type]
	if !exists {
		return fmt.Errorf("no handler registered for notification type: %s", notification.Type)
	}

	if err := handler.Handle(notification); err != nil {
		notification.Status = domain.StatusFailed
		notification.UpdatedAt = time.Now()
		if updateErr := s.repository.Update(notification); updateErr != nil {
			return fmt.Errorf("failed to update notification status: %w", updateErr)
		}
		return fmt.Errorf("failed to process notification: %w", err)
	}

	notification.Status = domain.StatusSent
	notification.UpdatedAt = time.Now()
	if err := s.repository.Update(notification); err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

func (s *NotificationService) CancelScheduledNotification(id string) error {
	notification, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if notification.ScheduledAt == nil || notification.ScheduledAt.Before(time.Now()) {
		return nil // Уведомление уже отправлено или не было запланировано
	}

	notification.Status = domain.StatusCancelled
	notification.UpdatedAt = time.Now()
	return s.repository.Update(notification)
}

func (s *NotificationService) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.stopCh:
			return nil
		case <-ticker.C:
			scheduledNotifications, err := s.repository.GetScheduledNotifications(time.Now())
			if err != nil {
				log.Printf("Error getting scheduled notifications: %v", err)
				continue
			}

			for _, notification := range scheduledNotifications {
				if err := s.processNotification(notification); err != nil {
					log.Printf("Error processing scheduled notification: %v", err)
				}
			}
		}
	}
}

func (s *NotificationService) Stop() {
	close(s.stopCh)
}

func (s *NotificationService) processNotification(notification *domain.Notification) error {
	s.mu.RLock()
	handler, exists := s.handlers[notification.Type]
	s.mu.RUnlock()

	if !exists {
		log.Printf("No handler registered for notification type: %s", notification.Type)
		return s.MarkAsFailed(notification, "no handler registered")
	}

	if err := handler.Handle(notification); err != nil {
		return s.MarkAsFailed(notification, err.Error())
	}

	return s.MarkAsSent(notification)
}

func (s *NotificationService) MarkAsSent(notification *domain.Notification) error {
	notification.Status = domain.StatusSent
	notification.UpdatedAt = time.Now()
	return s.repository.Update(notification)
}

func (s *NotificationService) MarkAsFailed(notification *domain.Notification, reason string) error {
	notification.Status = domain.StatusFailed
	notification.UpdatedAt = time.Now()
	return s.repository.Update(notification)
}

func (s *NotificationService) FindByID(id string) (*domain.Notification, error) {
	return s.repository.FindByID(id)
}
