package grpc

import (
	"context"
	"time"

	"github.com/XRS0/ToTalkB/notify/internal/domain/gen"
	"github.com/google/uuid"

	"github.com/XRS0/ToTalkB/notify/internal/domain"

	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedNotificationServiceServer
	service domain.NotificationService
}

func NewServer(service domain.NotificationService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) SendNotification(ctx context.Context, req *gen.SendNotificationRequest) (*gen.SendNotificationResponse, error) {
	notification := &domain.Notification{
		ID:      uuid.NewString(),
		Type:    req.Type,
		Payload: req.Payload,
	}

	if err := s.service.ProcessNotification(notification); err != nil {
		return nil, err
	}

	return &gen.SendNotificationResponse{
		Id:     notification.ID,
		Status: string(notification.Status),
	}, nil
}

func (s *Server) GetNotificationStatus(ctx context.Context, req *gen.GetNotificationStatusRequest) (*gen.GetNotificationStatusResponse, error) {
	notification, err := s.service.FindByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &gen.GetNotificationStatusResponse{
		Id:        notification.ID,
		Status:    string(notification.Status),
		CreatedAt: notification.CreatedAt.Format(time.RFC3339),
		UpdatedAt: notification.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func RegisterServer(s *grpc.Server, service domain.NotificationService) {
	gen.RegisterNotificationServiceServer(s, NewServer(service))
}
