package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/XRS0/ToTalkB/event_manager/config"
	gen "github.com/XRS0/ToTalkB/event_manager/internal/domain/gen"

	"github.com/XRS0/ToTalkB/event_manager/internal/domain/services"
	grpcImpl "github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/grpc"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/http/handlers"
	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/persistence/postgres"

	"github.com/XRS0/ToTalkB/event_manager/internal/infrastructure/notification"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	eventRepo := postgres.NewEventRepository(db)
	queueRepo := postgres.NewEventQueueRepository(db)

	// Initialize notification client
	notificationClient, err := notification.NewClient(fmt.Sprintf("%s:%d", cfg.NotificationService.Host, cfg.NotificationService.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to create notification client: %v", err)
	}

	// Initialize services
	eventService := services.NewEventService(eventRepo, notificationClient)
	queueService := services.NewEventQueueService(queueRepo)

	// Initialize gRPC server
	server := grpc.NewServer()
	eventServer := grpcImpl.NewEventServer(eventService)
	queueServer := grpcImpl.NewEventQueueServer(queueService)
	gen.RegisterEventServiceServer(server, eventServer)
	gen.RegisterEventQueueServiceServer(server, queueServer)

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.Server.GRPCPort)
		if err := server.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Initialize HTTP server
	router := gin.Default()
	eventHandler := handlers.NewEventHandler(eventService)
	router.GET("/api/events", eventHandler.GetEvents)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down servers...")
	server.GracefulStop()
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}
}
