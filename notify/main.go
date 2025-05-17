package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notify/internal/application"
	"notify/internal/config"
	"notify/internal/infrastructure/kafka"
	"notify/internal/infrastructure/persistence/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repository, err := postgres.NewRepository(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repository.Close()

	notificationService := application.NewNotificationApplicationService(repository)

	kafkaConfig := kafka.ConsumerConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topics: kafka.Topics{
			Notifications: cfg.Kafka.Topics.Notifications,
		},
	}
	consumer := kafka.NewConsumer(kafkaConfig, notificationService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
