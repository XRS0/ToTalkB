package websocket

import (
	"github.com/XRS0/ToTalkB/notify/internal/domain"
)

// NotificationHandler handles notifications via WebSocket
type NotificationHandler struct {
	manager *Manager
}

func NewNotificationHandler(manager *Manager) *NotificationHandler {
	return &NotificationHandler{
		manager: manager,
	}
}

func (h *NotificationHandler) Handle(notification *domain.Notification) error {
	return h.manager.SendToUser(notification.UserID, notification.Type, notification)
}
