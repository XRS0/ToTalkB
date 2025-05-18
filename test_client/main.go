package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gen "test_client/gen"
)

func main() {
	// Connect to event manager service
	eventConn, err := grpc.Dial("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to event manager: %v", err)
	}
	defer eventConn.Close()
	eventClient := gen.NewEventServiceClient(eventConn)
	queueClient := gen.NewEventQueueServiceClient(eventConn)

	// Connect to notification service
	notifyConn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to notification service: %v", err)
	}
	defer notifyConn.Close()
	notifyClient := gen.NewNotificationServiceClient(notifyConn)

	ctx := context.Background()

	// Test 1: Create an event
	payload := map[string]interface{}{
		"message": "Test event",
		"data": map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	eventResp, err := eventClient.ProcessEvent(ctx, &gen.ProcessEventRequest{
		Type:    "test_event",
		Source:  "test_client",
		Payload: payloadBytes,
	})
	if err != nil {
		log.Fatalf("Failed to process event: %v", err)
	}
	log.Printf("Event created: ID=%s, Status=%s", eventResp.Id, eventResp.Status)

	// Test 2: Get event status
	time.Sleep(time.Second) // Give some time for processing
	statusResp, err := eventClient.GetEventStatus(ctx, &gen.GetEventStatusRequest{
		Id: eventResp.Id,
	})
	if err != nil {
		log.Fatalf("Failed to get event status: %v", err)
	}
	log.Printf("Event status: ID=%s, Type=%s, Status=%s",
		statusResp.Id, statusResp.Type, statusResp.Status)

	// Test 3: Queue operations
	// Join queue
	joinResp, err := queueClient.JoinQueue(ctx, &gen.JoinQueueRequest{
		EventId: eventResp.Id,
		UserId:  "user1",
	})
	if err != nil {
		log.Fatalf("Failed to join queue: %v", err)
	}
	log.Printf("User joined queue: QueueID=%s, Position=%d", joinResp.QueueId, joinResp.Position)

	// Get queue status
	queueStatus, err := queueClient.GetQueueStatus(ctx, &gen.GetQueueStatusRequest{
		EventId: eventResp.Id,
	})
	if err != nil {
		log.Fatalf("Failed to get queue status: %v", err)
	}
	log.Printf("Queue status: %d users in queue", len(queueStatus.Queues))
	for _, q := range queueStatus.Queues {
		log.Printf("  - User: %s, Position: %d, Status: %s", q.UserId, q.Position, q.Status)
	}

	// Process next in queue
	nextResp, err := queueClient.ProcessNext(ctx, &gen.ProcessNextRequest{
		EventId: eventResp.Id,
	})
	if err != nil {
		log.Fatalf("Failed to process next: %v", err)
	}
	log.Printf("Processed next user: %s, status: %s", nextResp.Queue.UserId, nextResp.Queue.Status)

	// Test 4: Notification service tests
	// Send a direct notification
	notifyPayload := map[string]interface{}{
		"message":  "Test notification",
		"priority": "high",
	}
	notifyBytes, _ := json.Marshal(notifyPayload)

	notifyResp, err := notifyClient.SendNotification(ctx, &gen.SendNotificationRequest{
		Type:    "test_notification",
		Payload: notifyBytes,
	})
	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}
	log.Printf("Notification sent: ID=%s, Status=%s", notifyResp.Id, notifyResp.Status)

	// Get notification status
	time.Sleep(time.Second) // Give some time for processing
	notifyStatus, err := notifyClient.GetNotificationStatus(ctx, &gen.GetNotificationStatusRequest{
		Id: notifyResp.Id,
	})
	if err != nil {
		log.Fatalf("Failed to get notification status: %v", err)
	}
	log.Printf("Notification status: ID=%s, Status=%s, Created=%s, Updated=%s",
		notifyStatus.Id, notifyStatus.Status, notifyStatus.CreatedAt, notifyStatus.UpdatedAt)

	// Test 5: Send multiple notifications
	for i := 1; i <= 3; i++ {
		multiNotifyPayload := map[string]interface{}{
			"message":  fmt.Sprintf("Test notification %d", i),
			"priority": "medium",
		}
		multiNotifyBytes, _ := json.Marshal(multiNotifyPayload)

		multiNotifyResp, err := notifyClient.SendNotification(ctx, &gen.SendNotificationRequest{
			Type:    "batch_notification",
			Payload: multiNotifyBytes,
		})
		if err != nil {
			log.Printf("Failed to send batch notification %d: %v", i, err)
			continue
		}
		log.Printf("Batch notification %d sent: ID=%s", i, multiNotifyResp.Id)
	}
}
