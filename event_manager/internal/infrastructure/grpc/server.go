package grpc

import (
	"context"
	"time"

	"event_manager/internal/application"
	"event_manager/internal/domain"
	"event_manager/internal/domain/gen"

	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedEventServiceServer
	eventService *application.EventApplicationService
}

func NewServer(eventService *application.EventApplicationService) *Server {
	return &Server{
		eventService: eventService,
	}
}

func (s *Server) ProcessEvent(ctx context.Context, req *gen.ProcessEventRequest) (*gen.ProcessEventResponse, error) {
	event := &domain.Event{
		Type:    req.Type,
		Source:  req.Source,
		Payload: req.Payload,
	}

	if err := s.eventService.ProcessEvent(event); err != nil {
		return nil, err
	}

	return &gen.ProcessEventResponse{
		Id:     event.ID,
		Status: string(event.Status),
	}, nil
}

func (s *Server) GetEventStatus(ctx context.Context, req *gen.GetEventStatusRequest) (*gen.GetEventStatusResponse, error) {
	event, err := s.eventService.GetEventByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &gen.GetEventStatusResponse{
		Id:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Status:    string(event.Status),
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
		UpdatedAt: event.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func RegisterServer(s *grpc.Server, eventService *application.EventApplicationService) {
	gen.RegisterEventServiceServer(s, NewServer(eventService))
}
