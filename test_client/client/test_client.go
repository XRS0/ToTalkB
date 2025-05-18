package client

import (
	"context"
	"log"
	"time"

	pb "test_client/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient struct {
	eventClient      pb.EventServiceClient
	eventQueueClient pb.EventQueueServiceClient
	conn             *grpc.ClientConn
}

func NewTestClient(serverAddress string) (*TestClient, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &TestClient{
		eventClient:      pb.NewEventServiceClient(conn),
		eventQueueClient: pb.NewEventQueueServiceClient(conn),
		conn:             conn,
	}, nil
}

func (c *TestClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// EventService методы

func (c *TestClient) ProcessEvent(eventType, source string, payload []byte) (*pb.ProcessEventResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ProcessEventRequest{
		Type:    eventType,
		Source:  source,
		Payload: payload,
	}

	return c.eventClient.ProcessEvent(ctx, req)
}

func (c *TestClient) GetEventStatus(eventID string) (*pb.GetEventStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetEventStatusRequest{
		Id: eventID,
	}

	return c.eventClient.GetEventStatus(ctx, req)
}

// EventQueueService методы

func (c *TestClient) JoinQueue(eventID, userID string) (*pb.JoinQueueResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.JoinQueueRequest{
		EventId: eventID,
		UserId:  userID,
	}

	return c.eventQueueClient.JoinQueue(ctx, req)
}

func (c *TestClient) LeaveQueue(eventID, userID string) (*pb.LeaveQueueResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LeaveQueueRequest{
		EventId: eventID,
		UserId:  userID,
	}

	return c.eventQueueClient.LeaveQueue(ctx, req)
}

func (c *TestClient) GetQueueStatus(eventID string) (*pb.GetQueueStatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetQueueStatusRequest{
		EventId: eventID,
	}

	return c.eventQueueClient.GetQueueStatus(ctx, req)
}

func (c *TestClient) GetUserPosition(eventID, userID string) (*pb.GetUserPositionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetUserPositionRequest{
		EventId: eventID,
		UserId:  userID,
	}

	return c.eventQueueClient.GetUserPosition(ctx, req)
}

func (c *TestClient) ProcessNext(eventID string) (*pb.ProcessNextResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.ProcessNextRequest{
		EventId: eventID,
	}

	return c.eventQueueClient.ProcessNext(ctx, req)
}

func (c *TestClient) CloseQueue(eventID string) (*pb.CloseQueueResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CloseQueueRequest{
		EventId: eventID,
	}

	return c.eventQueueClient.CloseQueue(ctx, req)
}

// Пример использования
func Example() {
	client, err := NewTestClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Пример обработки события
	eventResp, err := client.ProcessEvent("test_event", "test_source", []byte(`{"key": "value"}`))
	if err != nil {
		log.Printf("Failed to process event: %v", err)
		return
	}
	log.Printf("Event processed: %s, status: %s", eventResp.Id, eventResp.Status)

	// Пример работы с очередью
	joinResp, err := client.JoinQueue(eventResp.Id, "user123")
	if err != nil {
		log.Printf("Failed to join queue: %v", err)
		return
	}
	log.Printf("Joined queue: %s, position: %d", joinResp.QueueId, joinResp.Position)

	// Получение статуса очереди
	queueStatus, err := client.GetQueueStatus(eventResp.Id)
	if err != nil {
		log.Printf("Failed to get queue status: %v", err)
		return
	}
	log.Printf("Queue status: %d entries", len(queueStatus.Queues))
}
