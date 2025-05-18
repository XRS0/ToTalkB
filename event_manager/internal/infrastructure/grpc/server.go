package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/application"
	gen "github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"

	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedEventServiceServer
	gen.UnimplementedEventQueueServiceServer
	eventService *application.EventApplicationService
}

func RegisterServer(grpcServer *grpc.Server, eventService *application.EventApplicationService) {
	server := &Server{
		eventService: eventService,
	}
	gen.RegisterEventServiceServer(grpcServer, server)
	gen.RegisterEventQueueServiceServer(grpcServer, server)
}

// EventService methods
func (s *Server) ProcessEvent(ctx context.Context, req *gen.ProcessEventRequest) (*gen.ProcessEventResponse, error) {
	event, err := s.eventService.ProcessEvent(ctx, req.Type, req.Source, req.Payload)
	if err != nil {
		return nil, err
	}

	return &gen.ProcessEventResponse{
		Id:     event.ID,
		Status: event.Status,
	}, nil
}

func (s *Server) GetEventStatus(ctx context.Context, req *gen.GetEventStatusRequest) (*gen.GetEventStatusResponse, error) {
	event, err := s.eventService.GetEventStatus(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &gen.GetEventStatusResponse{
		Id:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Status:    event.Status,
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
		UpdatedAt: event.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// EventQueueService methods
func (s *Server) JoinQueue(ctx context.Context, req *gen.JoinQueueRequest) (*gen.JoinQueueResponse, error) {
	queue, err := s.eventService.JoinQueue(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to join queue: %v", err)
	}

	return &gen.JoinQueueResponse{
		QueueId:  queue.ID,
		Position: int32(queue.Position),
	}, nil
}

func (s *Server) LeaveQueue(ctx context.Context, req *gen.LeaveQueueRequest) (*gen.LeaveQueueResponse, error) {
	success, err := s.eventService.LeaveQueue(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to leave queue: %v", err)
	}

	return &gen.LeaveQueueResponse{
		Success: success,
	}, nil
}

func (s *Server) GetQueueStatus(ctx context.Context, req *gen.GetQueueStatusRequest) (*gen.GetQueueStatusResponse, error) {
	queues, err := s.eventService.GetQueueStatus(ctx, req.EventId)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue status: %v", err)
	}

	response := &gen.GetQueueStatusResponse{
		Queues: make([]*gen.EventQueue, len(queues)),
	}

	for i, q := range queues {
		response.Queues[i] = &gen.EventQueue{
			Id:        q.ID,
			EventId:   q.EventID,
			UserId:    q.UserID,
			Status:    q.Status,
			Position:  int32(q.Position),
			CreatedAt: q.CreatedAt.Format(time.RFC3339),
			UpdatedAt: q.UpdatedAt.Format(time.RFC3339),
		}
	}

	return response, nil
}

func (s *Server) GetUserPosition(ctx context.Context, req *gen.GetUserPositionRequest) (*gen.GetUserPositionResponse, error) {
	position, err := s.eventService.GetUserPosition(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user position: %v", err)
	}

	return &gen.GetUserPositionResponse{
		Position: position,
	}, nil
}

func (s *Server) ProcessNext(ctx context.Context, req *gen.ProcessNextRequest) (*gen.ProcessNextResponse, error) {
	queue, err := s.eventService.ProcessNext(ctx, req.EventId)
	if err != nil {
		return nil, fmt.Errorf("failed to process next: %v", err)
	}

	return &gen.ProcessNextResponse{
		Queue: &gen.EventQueue{
			Id:        queue.ID,
			EventId:   queue.EventID,
			UserId:    queue.UserID,
			Status:    queue.Status,
			Position:  int32(queue.Position),
			CreatedAt: queue.CreatedAt.Format(time.RFC3339),
			UpdatedAt: queue.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *Server) CloseQueue(ctx context.Context, req *gen.CloseQueueRequest) (*gen.CloseQueueResponse, error) {
	success, err := s.eventService.CloseQueue(ctx, req.EventId)
	if err != nil {
		return nil, fmt.Errorf("failed to close queue: %v", err)
	}

	return &gen.CloseQueueResponse{
		Success: success,
	}, nil
}
