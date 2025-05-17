package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"event_manager/internal/application"
	"event_manager/internal/config"
	grpcserver "event_manager/internal/infrastructure/grpc"
	"event_manager/internal/infrastructure/persistence/postgres"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName,
	)
	repository, err := postgres.NewRepository(connStr)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repository.Close()

	eventService := application.NewEventApplicationService(repository)

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
