package main

import (
	"encoding/json"
	"log"
	"test_client/client"
	"time"
)

func main() {
	// Создаем клиент
	client, err := client.NewTestClient("localhost:9091")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Тестируем EventService
	testEventService(client)

	// Тестируем EventQueueService
	testEventQueueService(client)
}

func testEventService(client *client.TestClient) {
	log.Println("Testing EventService...")

	// Создаем тестовое событие
	payload := map[string]interface{}{
		"message": "Test event",
		"data": map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}
	payloadBytes, _ := json.Marshal(payload)

	// Обрабатываем событие
	eventResp, err := client.ProcessEvent("test_event", "test_client", payloadBytes)
	if err != nil {
		log.Printf("Failed to process event: %v", err)
		return
	}
	log.Printf("Event processed successfully: ID=%s, Status=%s", eventResp.Id, eventResp.Status)

	// Получаем статус события
	time.Sleep(time.Second) // Даем время на обработку
	statusResp, err := client.GetEventStatus(eventResp.Id)
	if err != nil {
		log.Printf("Failed to get event status: %v", err)
		return
	}
	log.Printf("Event status: ID=%s, Type=%s, Status=%s",
		statusResp.Id, statusResp.Type, statusResp.Status)
}

func testEventQueueService(client *client.TestClient) {
	log.Println("\nTesting EventQueueService...")

	// Создаем тестовое событие для очереди
	eventResp, err := client.ProcessEvent("queue_test", "test_client", []byte(`{"type": "queue_test"}`))
	if err != nil {
		log.Printf("Failed to create test event: %v", err)
		return
	}
	eventID := eventResp.Id

	// Добавляем нескольких пользователей в очередь
	users := []string{"user1", "user2", "user3"}
	for _, userID := range users {
		joinResp, err := client.JoinQueue(eventID, userID)
		if err != nil {
			log.Printf("Failed to join queue for user %s: %v", userID, err)
			continue
		}
		log.Printf("User %s joined queue: position=%d", userID, joinResp.Position)
	}

	// Получаем статус очереди
	queueStatus, err := client.GetQueueStatus(eventID)
	if err != nil {
		log.Printf("Failed to get queue status: %v", err)
	} else {
		log.Printf("Queue status: %d entries", len(queueStatus.Queues))
		for _, q := range queueStatus.Queues {
			log.Printf("  - User: %s, Position: %d, Status: %s", q.UserId, q.Position, q.Status)
		}
	}

	// Проверяем позицию конкретного пользователя
	position, err := client.GetUserPosition(eventID, "user2")
	if err != nil {
		log.Printf("Failed to get user position: %v", err)
	} else {
		log.Printf("User2 position: %d", position.Position)
	}

	// Обрабатываем следующего в очереди
	nextResp, err := client.ProcessNext(eventID)
	if err != nil {
		log.Printf("Failed to process next: %v", err)
	} else {
		log.Printf("Processed next user: %s, status: %s", nextResp.Queue.UserId, nextResp.Queue.Status)
	}

	// Удаляем пользователя из очереди
	leaveResp, err := client.LeaveQueue(eventID, "user3")
	if err != nil {
		log.Printf("Failed to leave queue: %v", err)
	} else {
		log.Printf("User left queue: success=%v", leaveResp.Success)
	}

	// Закрываем очередь
	closeResp, err := client.CloseQueue(eventID)
	if err != nil {
		log.Printf("Failed to close queue: %v", err)
	} else {
		log.Printf("Queue closed: success=%v", closeResp.Success)
	}

	// Финальный статус очереди
	finalStatus, err := client.GetQueueStatus(eventID)
	if err != nil {
		log.Printf("Failed to get final queue status: %v", err)
	} else {
		log.Printf("Final queue status: %d entries", len(finalStatus.Queues))
		for _, q := range finalStatus.Queues {
			log.Printf("  - User: %s, Position: %d, Status: %s", q.UserId, q.Position, q.Status)
		}
	}
}
