package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/XRS0/ToTalkB/event_manager/internal/application"
	"github.com/XRS0/ToTalkB/event_manager/internal/config"
	grpcserver "github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/grpc"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/notification"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/persistence/memory"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize notification client
	notificationClient, err := notification.NewClient(&cfg.NotificationService)
	if err != nil {
		log.Fatalf("Failed to create notification client: %v", err)
	}
	defer notificationClient.Close()

	repository := memory.NewInMemoryEventRepository()

	// Pass notification client to the application service
	eventService := application.NewEventApplicationService(repository, notificationClient)

	grpcServer := grpc.NewServer()
	grpcserver.RegisterServer(grpcServer, eventService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
}
