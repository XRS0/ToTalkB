package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notify/internal/config"
	"notify/internal/infrastructure/persistence/postgres"
	"notify/internal/server"
	"notify/internal/service"
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

	notificationService := service.NewNotificationService()
	srv := server.NewServer(cfg, notificationService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
