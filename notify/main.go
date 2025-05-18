package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/XRS0/ToTalkB/notify/internal/handlers"
	"github.com/XRS0/ToTalkB/notify/internal/infrastructure/persistence/postgres"
	"github.com/XRS0/ToTalkB/notify/internal/server"
	"github.com/XRS0/ToTalkB/notify/internal/service"

	"github.com/XRS0/ToTalkB/notify/internal/config"
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

	notificationService := service.NewNotificationService(repository)

	// Регистрируем обработчики уведомлений
	notificationService.RegisterHandler("push", handlers.NewPushHandler())

	srv := server.NewServer(cfg, notificationService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем сервис уведомлений в отдельной горутине
	go func() {
		if err := notificationService.Start(ctx); err != nil {
			log.Printf("Notification service error: %v", err)
		}
	}()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	notificationService.Stop()
}
