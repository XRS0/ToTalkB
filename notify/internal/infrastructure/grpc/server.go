package grpc

import (
	"context"
	"time"

	"notify/internal/domain"
	"notify/internal/domain/gen"

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
	// Implementation would depend on your repository interface
	// This is a placeholder
	return &gen.GetNotificationStatusResponse{
		Id:        req.Id,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func RegisterServer(s *grpc.Server, service domain.NotificationService) {
	gen.RegisterNotificationServiceServer(s, NewServer(service))
}
