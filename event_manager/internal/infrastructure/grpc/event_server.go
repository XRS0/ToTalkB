package grpc

import (
	"context"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"
	"github.com/XRS0/ToTalkB/event_manager/internal/domain/services"
)

type EventServer struct {
	gen.UnimplementedEventServiceServer
	service *services.EventService
}

func NewEventServer(service *services.EventService) *EventServer {
	return &EventServer{service: service}
}

func (s *EventServer) GetAllEvents(ctx context.Context, req *gen.GetAllEventsRequest) (*gen.GetAllEventsResponse, error) {
	events, err := s.service.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}
	return &gen.GetAllEventsResponse{Events: events}, nil
}
