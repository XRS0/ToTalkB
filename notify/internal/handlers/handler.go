package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/XRS0/ToTalkB/notify/internal/domain"
)

// BaseHandler - базовый обработчик уведомлений
type BaseHandler struct {
	Type string
}

// PushNotification - структура для push уведомлений
type PushNotification struct {
	DeviceToken string `json:"device_token"`
	Title       string `json:"title"`
	Body        string `json:"body"`
}

// PushHandler - обработчик push уведомлений
type PushHandler struct {
	BaseHandler
}

func NewPushHandler() *PushHandler {
	return &PushHandler{
		BaseHandler: BaseHandler{
			Type: "push",
		},
	}
}

func (h *PushHandler) Handle(notification *domain.Notification) error {
	var pushData PushNotification
	if err := json.Unmarshal(notification.Payload, &pushData); err != nil {
		return fmt.Errorf("failed to unmarshal push notification: %w", err)
	}

	// TODO: Здесь должна быть реальная отправка push уведомления
	log.Printf("Sending push notification to device %s: %s", pushData.DeviceToken, pushData.Title)
	return nil
}
