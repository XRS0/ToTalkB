package grpc

import (
	"context"
	"time"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain"
	"github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"
	"github.com/XRS0/ToTalkB/event_manager/internal/domain/services"
)

type EventQueueServer struct {
	gen.UnimplementedEventQueueServiceServer
	service *services.EventQueueService
}

func NewEventQueueServer(service *services.EventQueueService) *EventQueueServer {
	return &EventQueueServer{service: service}
}

func toProtoEventQueue(queue *domain.EventQueue) *gen.EventQueue {
	return &gen.EventQueue{
		Id:        queue.ID,
		EventId:   queue.EventID,
		UserId:    queue.UserID,
		Status:    queue.Status,
		Position:  int32(queue.Position),
		CreatedAt: queue.CreatedAt.Format(time.RFC3339),
		UpdatedAt: queue.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *EventQueueServer) GetQueueStatus(ctx context.Context, req *gen.GetQueueStatusRequest) (*gen.GetQueueStatusResponse, error) {
	queues, err := s.service.GetQueueStatus(ctx, req.EventId)
	if err != nil {
		return nil, err
	}

	protoQueues := make([]*gen.EventQueue, len(queues))
	for i, q := range queues {
		protoQueues[i] = toProtoEventQueue(q)
	}
	return &gen.GetQueueStatusResponse{Queues: protoQueues}, nil
}

func (s *EventQueueServer) GetUserPosition(ctx context.Context, req *gen.GetUserPositionRequest) (*gen.GetUserPositionResponse, error) {
	position, err := s.service.GetUserPosition(ctx, req.EventId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &gen.GetUserPositionResponse{Position: int32(position)}, nil
}

func (s *EventQueueServer) ProcessNext(ctx context.Context, req *gen.ProcessNextRequest) (*gen.ProcessNextResponse, error) {
	queue, err := s.service.ProcessNext(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return &gen.ProcessNextResponse{Queue: toProtoEventQueue(queue)}, nil
}
