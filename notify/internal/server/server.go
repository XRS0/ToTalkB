package server

import (
	"context"
	"fmt"
	"log"
	"net"

	grpcserver "github.com/XRS0/ToTalkB/notify/internal/infrastructure/grpc"
	"github.com/XRS0/ToTalkB/notify/internal/service"

	"github.com/XRS0/ToTalkB/notify/internal/websocket"

	"github.com/XRS0/ToTalkB/notify/internal/config"

	"google.golang.org/grpc"
)

type Server struct {
	config     *config.Config
	httpServer *HTTPServer
	grpcServer *grpc.Server
	wsManager  *websocket.Manager
	service    *service.NotificationService
}

func NewServer(cfg *config.Config, svc *service.NotificationService) *Server {
	// Создаем WebSocket менеджер
	wsManager := websocket.NewManager()
	go wsManager.Start()

	// Создаем WebSocket обработчик уведомлений
	wsHandler := websocket.NewNotificationHandler(wsManager)
	svc.RegisterHandler("websocket", wsHandler)

	// Создаем HTTP сервер с JWT ключом
	httpServer := NewHTTPServer(wsManager, []byte(cfg.Auth.JWTKey))

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	grpcserver.RegisterServer(grpcServer, svc)

	return &Server{
		config:     cfg,
		service:    svc,
		grpcServer: grpcServer,
		httpServer: httpServer,
		wsManager:  wsManager,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Запускаем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.GRPCPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Starting gRPC server on %s", fmt.Sprintf(":%d", s.config.Server.GRPCPort))
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Запускаем HTTP сервер
	return s.httpServer.Start(fmt.Sprintf(":%d", s.config.Server.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.httpServer.Shutdown(ctx)
}
