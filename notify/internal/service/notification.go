package service

import (
	"context"
	"log"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/domain"
)

type NotificationService struct {
	handlers map[string]NotificationHandler
}

type NotificationHandler interface {
	Handle(notification *domain.Notification) error
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		handlers: make(map[string]NotificationHandler),
	}
}

func (s *NotificationService) RegisterHandler(notificationType string, handler NotificationHandler) {
	s.handlers[notificationType] = handler
}

func (s *NotificationService) ProcessNotification(notification *domain.Notification) error {
	notification.Status = domain.StatusPending
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	handler, exists := s.handlers[notification.Type]
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
	return nil
}

func (s *NotificationService) MarkAsFailed(notification *domain.Notification, reason string) error {
	notification.Status = domain.StatusFailed
	notification.UpdatedAt = time.Now()
	return nil
}

func (s *NotificationService) Start(ctx context.Context) error {
	// Service is started by the Kafka consumer
	return nil
}
