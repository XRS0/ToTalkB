package grpc

import (
	"context"
	"time"

	"event_manager/internal/domain"
	pb "event_manager/internal/domain/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventServer struct {
	pb.UnimplementedEventServiceServer
	pb.UnimplementedEventQueueServiceServer
	eventService domain.EventService
	queueService domain.EventQueueService
}

func NewEventServer(eventService domain.EventService, queueService domain.EventQueueService) *EventServer {
	return &EventServer{
		eventService: eventService,
		queueService: queueService,
	}
}

// EventService реализация

func (s *EventServer) ProcessEvent(ctx context.Context, req *pb.ProcessEventRequest) (*pb.ProcessEventResponse, error) {
	event := &domain.Event{
		Type:    req.Type,
		Source:  req.Source,
		Payload: req.Payload,
	}

	err := s.eventService.ProcessEvent(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProcessEventResponse{
		Id:     event.ID,
		Status: string(event.Status),
	}, nil
}

func (s *EventServer) GetEventStatus(ctx context.Context, req *pb.GetEventStatusRequest) (*pb.GetEventStatusResponse, error) {
	event, err := s.eventService.GetEventStatus(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEventStatusResponse{
		Id:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Status:    string(event.Status),
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
		UpdatedAt: event.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// EventQueueService реализация

func (s *EventServer) JoinQueue(ctx context.Context, req *pb.JoinQueueRequest) (*pb.JoinQueueResponse, error) {
	err := s.queueService.JoinQueue(req.EventId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	position, err := s.queueService.GetUserPosition(req.EventId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.JoinQueueResponse{
		QueueId:  req.EventId, // Используем event_id как queue_id
		Position: int32(position),
	}, nil
}

func (s *EventServer) LeaveQueue(ctx context.Context, req *pb.LeaveQueueRequest) (*pb.LeaveQueueResponse, error) {
	err := s.queueService.LeaveQueue(req.EventId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LeaveQueueResponse{
		Success: true,
	}, nil
}

func (s *EventServer) GetQueueStatus(ctx context.Context, req *pb.GetQueueStatusRequest) (*pb.GetQueueStatusResponse, error) {
	queues, err := s.queueService.GetQueueStatus(req.EventId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbQueues := make([]*pb.EventQueue, len(queues))
	for i, q := range queues {
		pbQueues[i] = &pb.EventQueue{
			Id:        q.ID,
			EventId:   q.EventID,
			UserId:    q.UserID,
			Status:    q.Status,
			Position:  int32(q.Position),
			CreatedAt: q.CreatedAt.Format(time.RFC3339),
			UpdatedAt: q.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &pb.GetQueueStatusResponse{
		Queues: pbQueues,
	}, nil
}

func (s *EventServer) GetUserPosition(ctx context.Context, req *pb.GetUserPositionRequest) (*pb.GetUserPositionResponse, error) {
	position, err := s.queueService.GetUserPosition(req.EventId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserPositionResponse{
		Position: int32(position),
	}, nil
}

func (s *EventServer) ProcessNext(ctx context.Context, req *pb.ProcessNextRequest) (*pb.ProcessNextResponse, error) {
	err := s.queueService.ProcessNext(req.EventId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Получаем обновленный статус очереди
	queues, err := s.queueService.GetQueueStatus(req.EventId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Находим активную запись
	var activeQueue *domain.EventQueue
	for _, q := range queues {
		if q.Status == string(domain.QueueStatusActive) {
			activeQueue = q
			break
		}
	}

	if activeQueue == nil {
		return nil, status.Error(codes.NotFound, "no active queue entry found")
	}

	return &pb.ProcessNextResponse{
		Queue: &pb.EventQueue{
			Id:        activeQueue.ID,
			EventId:   activeQueue.EventID,
			UserId:    activeQueue.UserID,
			Status:    activeQueue.Status,
			Position:  int32(activeQueue.Position),
			CreatedAt: activeQueue.CreatedAt.Format(time.RFC3339),
			UpdatedAt: activeQueue.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *EventServer) CloseQueue(ctx context.Context, req *pb.CloseQueueRequest) (*pb.CloseQueueResponse, error) {
	err := s.queueService.CloseQueue(req.EventId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CloseQueueResponse{
		Success: true,
	}, nil
}
