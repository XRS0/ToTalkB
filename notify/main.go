package main

import (
	"context"
	"fmt"
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

	// Initialize repository
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	repository, err := postgres.NewRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Initialize application service
	notificationService := application.NewNotificationApplicationService(repository)

	// Initialize Kafka consumer
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

	// Start Kafka consumer
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
