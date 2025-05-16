package application

import (
	"log"
	"time"

	"notify/internal/domain"
)

type NotificationApplicationService struct {
	repository domain.NotificationRepository
	handlers   map[string]NotificationHandler
}

type NotificationHandler interface {
	Handle(notification *domain.Notification) error
}

func NewNotificationApplicationService(repository domain.NotificationRepository) *NotificationApplicationService {
	return &NotificationApplicationService{
		repository: repository,
		handlers:   make(map[string]NotificationHandler),
	}
}

func (s *NotificationApplicationService) RegisterHandler(notificationType string, handler NotificationHandler) {
	s.handlers[notificationType] = handler
}

func (s *NotificationApplicationService) ProcessNotification(notification *domain.Notification) error {
	notification.Status = domain.StatusPending
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	if err := s.repository.Save(notification); err != nil {
		return err
	}

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

func (s *NotificationApplicationService) MarkAsSent(notification *domain.Notification) error {
	notification.Status = domain.StatusSent
	notification.UpdatedAt = time.Now()
	return s.repository.Update(notification)
}

func (s *NotificationApplicationService) MarkAsFailed(notification *domain.Notification, reason string) error {
	notification.Status = domain.StatusFailed
	notification.UpdatedAt = time.Now()
	return s.repository.Update(notification)
}
